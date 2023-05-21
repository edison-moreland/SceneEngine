use "collections"
use "debug"
use "buffered"

type MsgId is U32
type SendMsg is {(MsgId, (Array[U8 val] iso | None))} val
type ReceiveMsg is {(MsgId, Array[U8 val] iso)} val

primitive StartSubMsg
    fun apply(stdin: InputStream, stdout: OutStream, receiver: ReceiveMsg): SendMsg =>
        stdin(_MsgReceiver(receiver))
        _MsgSender(stdout)~apply()

actor _MsgSender
    let out: OutStream

    new create(out': OutStream) =>
        out = out'

    be apply(id: MsgId, data: (Array[U8] iso | None)) =>
        // TODO: Rewrite this using https://stdlib.ponylang.io/buffered-Writer/

        let data_size: USize = match data
        | None => 0
        | let d: Array[U8] iso => d.size()
        end

        let msg_header: Array[U8] iso = recover Array[U8](data_size+8) end
        msg_header.push_u32(id)
        msg_header.push_u32(data_size.u32())

        out.write(consume msg_header)
        match data
        | let d: Array[U8] iso => out.write(consume d)
        end
        out.flush()

class _MsgReceiver is InputNotify
    let _receive_msg: ReceiveMsg

    new iso create(receive_msg: ReceiveMsg) =>
        _receive_msg = receive_msg

    fun ref apply(data: Array[U8 val] iso) =>
        let rb = Reader.>append(consume data)

        try
            let msg_id = rb.u32_le()?
            let msg_len = rb.u32_le()?

            _receive_msg(msg_id, rb.block(msg_len.usize())?)
        else
            // TODO: handle this. What happens if the chunk we get is too small?
            return
        end
