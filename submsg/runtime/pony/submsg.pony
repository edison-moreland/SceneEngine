use "collections"
use "debug"
use "buffered"

type MsgId is U32
type SendMsg is {(MsgId, (Array[U8 val] iso | None))} val
type ReceiveMsg is {(MsgId, Array[U8 val] iso)} val

primitive StartSubMsg
    fun apply(stdin: InputStream, stdout: OutStream, receiver: ReceiveMsg): SendMsg =>
        stdin(_MsgReceiver(receiver) where chunk_size = 512)
        _MsgSender(stdout)~apply()

actor _MsgSender
    let buf: Writer = Writer
    let out: OutStream

    new create(out': OutStream) =>
        out = out'

    be apply(id: MsgId, data: (Array[U8] iso | None)) =>
        buf.u32_le(id)

        match data
        | None =>
            buf.u32_le(0)
        | let d: Array[U8] iso =>
            buf.u32_le(d.size().u32())
            buf.write(consume d)
        end

        for chunk in buf.done().values() do
            out.write(chunk)
        end
        out.flush()

class _MsgReceiver is InputNotify
    let _receive_msg: ReceiveMsg
    let _buf: Reader = Reader

    new iso create(receive_msg: ReceiveMsg) =>
        _receive_msg = receive_msg

    fun ref apply(data: Array[U8 val] iso) =>
        _buf.append(consume data)

        let msg_len' = try _buf.peek_u32_le(where offset=4)? else 8 end
        if msg_len'.usize() > _buf.size() then
            return // This message is coming in more than one chunk
        end

        (let msg_id, let msg_len) = try
            (_buf.u32_le()?, _buf.u32_le()?)
        else
            Debug("Error receiving header!" where stream=DebugErr)
            return
        end

        try
            _receive_msg(msg_id, _buf.block(msg_len.usize())?)
        else
            Debug("Error receiving header!" where stream=DebugErr)
        end