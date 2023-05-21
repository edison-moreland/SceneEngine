use "msgpack"
use "buffered"
use "debug"

primitive Collect
    fun apply(chunks: Array[(String val | Array[U8 val] val)] box): Array[U8 val] iso^ =>
        let b = recover iso Array[U8] end
        for chunk in chunks.values() do
            b.append(chunk)
        end

        consume b

class val MsgCoreInfo
    let version: String

    new create(version': String) =>
        version = version'

    fun marshal(): Array[U8 val] iso^ =>
        var writer: Writer ref = Writer

        try
            MessagePackEncoder.fixmap(writer, 1)?
            MessagePackEncoder.str_32(writer, "Version")?
            MessagePackEncoder.str_32(writer, version)?
        else
            Debug("error in MsgCoreInfo.marshal" where stream = DebugErr)
        end

        Collect(writer.done())
