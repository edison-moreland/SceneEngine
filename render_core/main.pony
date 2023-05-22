use "messages"
use "format"
use "random"
use "collections"
use "../submsg/runtime/pony"

actor Main is CoreServer
    let client: EngineClient
    let env: Env

    var rand: Rand = Rand

    var render_config: Config = Config(0, 0)

    new create(env': Env) =>
        env = env'
        env.err.print("PonyCore: starting")

        let sendMsg = StartSubMsg(env.input, env.out, CoreRouter(this))
        client = EngineClient(sendMsg)

        client.core_ready(None)
        env.err.print("PonyCore: ready")

    be info(body: Array[U8] iso) =>
        client.core_info(Marshal(MsgCoreInfo("PonyCore v0.0.1")))

    be config(body: Array[U8] iso) =>
        render_config = Config.unmarshal_msgpack(consume body)
        env.err.print("PonyCore: render_config x" +
            Format.int[U64](render_config.image_width) + " y" +
            Format.int[U64](render_config.image_height))

        client.core_ready(None)

    be render_frame(body: Array[U8] iso) =>
        let pixel_buf = Array[MsgPackMarshalable](render_config.image_width.usize())

        for y in Range(0, render_config.image_height.usize()) do
            for x in Range(0, render_config.image_width.usize()) do
                pixel_buf.push(Pixel(x.u64(), y.u64(), Color(rand.real(), rand.real(), rand.real())))
            end
            client.pixel_batch(Marshal(pixel_buf))
            pixel_buf.clear()
        end

        client.core_ready(None)
