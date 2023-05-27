use "format"
use "random"
use "collections"
use "buffered"
use "runtime_info"

use "logger"
use "renderer"
use "messages"
use scene = "scene"
use "../submsg/runtime/pony"

actor Main is CoreServer
    let env: Env
    let client: EngineClient
    let logger: Logger

    var rand: Rand = Rand

    var render_config: Config = Config(where
            aspect_ratio' = 0.0,
            image_width' = 0,
            image_height' = 0,
            samples' = 0,
            depth' = 0
        )

    new create(env': Env) =>
        env = env'

        logger = Logger(env.err)
        logger.log("PonyCore: starting")

        let sendMsg = StartSubMsg(env.input, env.out, CoreRouter(this))
        client = EngineClient(sendMsg)

        client.core_ready(None)
        logger.log("PonyCore: ready")

    be info(body: Array[U8] iso) =>
        client.core_info(Marshal(MsgCoreInfo("PonyCore v0.0.1")))

    be config(body: Array[U8] iso) =>
        let r: Reader = Reader
        r.append(consume body)

        render_config = UnmarshalMsgPackConfig(consume r)
        logger.log("PonyCore: render_config x" +
            Format.int[U64](render_config.image_width) + " y" +
            Format.int[U64](render_config.image_height))

        client.core_ready(None)

    be render_frame(body: Array[U8] iso) =>
        let r: Reader = Reader
        r.append(consume body)

        let frame_scene = scene.Transform(UnmarshalMsgPackScene(consume r))

        let pixel = PixelBatcher(client, 100) // TODO: Best pixel batch size?

        Renderer.render(
            SchedulerInfoAuth(env.root),
            render_config,
            frame_scene,
            pixel~apply(),
            {()(client) =>
                pixel.send_batch()
                pixel.sync({() =>
                    client.core_ready(None)
                })
            }
        )

actor PixelBatcher
    let client: EngineClient

    let buffer: Array[MsgPackMarshalable]
    let batch_size: USize

    new create(client': EngineClient, batch_size': USize) =>
        client = client'

        batch_size = batch_size'
        buffer = Array[MsgPackMarshalable](batch_size)

    be apply(x: U64, y: U64, color: PixelColor) =>
        buffer.push(Pixel(where
            x' = x,
            y' = y,
            color' = Color(where
                r' = color._1,
                g' = color._2,
                b' = color._3
            )
        ))

        if buffer.size() >= batch_size then
            send_batch()
        end

    be send_batch() =>
        client.pixel_batch(Marshal(buffer))
        buffer.clear()

    be sync(f: {()} val) =>
        f()