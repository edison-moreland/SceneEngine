// Code generated by submsg; DO NOT EDIT. 
use "../../submsg/runtime/pony" 

primitive Engine
    fun core_ready(): MsgId => 0
    fun core_info(): MsgId => 1
    fun pixel_batch(): MsgId => 2

interface tag EngineServer
    fun tag core_ready(body: Array[U8] iso)
    fun tag core_info(body: Array[U8] iso)
    fun tag pixel_batch(body: Array[U8] iso)

primitive EngineRouter
    fun apply(s: EngineServer): ReceiveMsg =>
        {(id: MsgId, body: Array[U8] iso) =>
            match id
            | Engine.core_ready() => s.core_ready(consume body)
            | Engine.core_info() => s.core_info(consume body)
            | Engine.pixel_batch() => s.pixel_batch(consume body)
            end
        }

actor EngineClient
    let send_msg: SendMsg

    new create(send_msg': SendMsg) =>
        send_msg = send_msg'

    be core_ready(data: (Array[U8 val] iso | None)) =>
        send_msg(Engine.core_ready(), consume data)

    be core_info(data: (Array[U8 val] iso | None)) =>
        send_msg(Engine.core_info(), consume data)

    be pixel_batch(data: (Array[U8 val] iso | None)) =>
        send_msg(Engine.pixel_batch(), consume data)


primitive Core
    fun info(): MsgId => 0
    fun render_frame(): MsgId => 1

interface tag CoreServer
    fun tag info(body: Array[U8] iso)
    fun tag render_frame(body: Array[U8] iso)

primitive CoreRouter
    fun apply(s: CoreServer): ReceiveMsg =>
        {(id: MsgId, body: Array[U8] iso) =>
            match id
            | Core.info() => s.info(consume body)
            | Core.render_frame() => s.render_frame(consume body)
            end
        }

actor CoreClient
    let send_msg: SendMsg

    new create(send_msg': SendMsg) =>
        send_msg = send_msg'

    be info(data: (Array[U8 val] iso | None)) =>
        send_msg(Core.info(), consume data)

    be render_frame(data: (Array[U8 val] iso | None)) =>
        send_msg(Core.render_frame(), consume data)

