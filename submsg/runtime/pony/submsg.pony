use "collections"
use "debug"

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
    var _recv_buf: Array[U8] ref = Array[U8]()
    var _targetSize: USize = 8

    let _receive_msg: ReceiveMsg

    new iso create(receive_msg: ReceiveMsg) =>
        _receive_msg = receive_msg

    fun ref apply(data: Array[U8 val] iso) =>
        // We could be sent data in any size chunk
        // We need to collect the entire message into _recv_buf before sending it off

        _recv_buf.append(consume data)

        if _recv_buf.size() < _targetSize then
            return
        end

        (let msg_id, let msg_length) = try
            (_recv_buf.read_u32(0)?, _recv_buf.read_u32(4)?)
        else
            return
        end

        if _recv_buf.size() < (msg_length+8).usize() then
            _targetSize = (msg_length+8).usize()
            return
        end

        // Buffer is full of goodies!
        let msg_buf: Array[U8] iso = recover Array[U8](msg_length.usize()) end
        for v in _recv_buf.slice(where from=8).values() do
            msg_buf.push(v)
        end

        _receive_msg(msg_id, consume msg_buf)

        _recv_buf.clear()
        _targetSize = 8