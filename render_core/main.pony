use "messages"
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
        client.coreInfo(MsgCoreInfo("PonyCore v0.0.1").marshal())