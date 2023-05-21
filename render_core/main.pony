use "messages"
use "../submsg/runtime/pony"

actor Main is CoreServer
    let client: EngineClient

    new create(env: Env) =>
        env.err.print("PonyCore: starting")

        let sendMsg = StartSubMsg(env.input, env.out, CoreRouter(this))
        client = EngineClient(sendMsg)

        client.core_ready(None)
        env.err.print("PonyCore: ready")

    be info(body: Array[U8] iso) =>
        client.core_info(Marshal(MsgCoreInfo("PonyCore v0.0.1")))

    be render_frame(body: Array[U8] iso) =>
        let pixel = Array[MsgPackMarshalable](3)

        pixel.push(Pixel(1, 1, Color(0.4, 0.1, 0.4)))
        pixel.push(Pixel(1, 2, Color(0.2, 0.5, 0.8)))
        pixel.push(Pixel(1, 3, Color(0.5, 0.0, 0.1)))

        client.pixel_batch(Marshal(pixel))

        pixel.clear()
        pixel.push(Pixel(2, 1, Color(0.2, 0.1, 0.8)))
        pixel.push(Pixel(2, 2, Color(0.1, 0.23, 0.1)))
        pixel.push(Pixel(2, 3, Color(0.4, 0.12, 0.9)))

        client.pixel_batch(Marshal(pixel))
        client.core_ready(None)

