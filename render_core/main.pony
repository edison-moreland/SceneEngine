use "../submsg/runtime/pony"

actor Main
    let client: EngineClient

    new create(env: Env) =>
        env.err.print("PonyCore: starting")

        let sendMsg = StartSubMsg(env.input, env.out, CoreRouter(this))
        client = EngineClient(sendMsg)

        client.coreReady(None)
        env.err.print("PonyCore: ready")

    be info(body: Array[U8] iso) =>
        let i = recover String end
        i.append("PonyCore v0.0.1")

        client.coreInfo((consume i).iso_array())