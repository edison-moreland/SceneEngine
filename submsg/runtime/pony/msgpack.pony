use "buffered"
use "msgpack"
use "debug"
use "collections"

interface val MsgPackMarshalable
    fun marshal_msgpack(b: Writer ref)?

primitive Marshal
    fun apply(m: (MsgPackMarshalable | Array[MsgPackMarshalable])): Array[U8 val] iso^ =>
        var writer: Writer ref = Writer

        try
            match m
            | let one: MsgPackMarshalable =>
                one.marshal_msgpack(writer)?
            | let many: Array[MsgPackMarshalable] =>
                array_header(writer, many.size())?
                for m' in many.values() do
                    m'.marshal_msgpack(writer)?
                end
            end
        else
            Debug("error marshalling!" where stream = DebugErr)
        end

        collect(writer.done())

    fun array_header(b: Writer ref, size: USize)? =>
        // Write a message pack header for the smallest array that can hold size items
        if size <= 15 then
            MessagePackEncoder.fixarray(b, size.u8())?
        elseif size <= 65535 then
            MessagePackEncoder.array_16(b, size.u16())
        else
            MessagePackEncoder.array_32(b, size.u32())
        end


    fun collect(chunks: Array[(String val | Array[U8 val] val)] box): Array[U8 val] iso^ =>
        let b = recover iso Array[U8] end
        for chunk in chunks.values() do
            b.append(chunk)
        end

        consume b

primitive Unmarshal
    fun map(b: Reader ref): USize val ? =>
        try
            MessagePackDecoder.fixmap(b)?.usize()
        else
            try
                MessagePackDecoder.map_16(b)?.usize()
            else
                MessagePackDecoder.map_32(b)?.usize()
            end
        end

    fun array_header(b: Reader ref): USize val ? =>
        try
            MessagePackDecoder.fixarray(b)?.usize()
        else
            try
                MessagePackDecoder.array_16(b)?.usize()
            else
                MessagePackDecoder.array_32(b)?.usize()
            end
        end

    fun array[M: MsgPackMarshalable #send](b: Reader ref, u: {(Reader ref): M^?}): Array[M] iso^ ? =>
        let array_size = array_header(b)?
        let res: Array[M] iso = recover Array[M](array_size) end

        for _ in Range(0, array_size) do
            res.push(u(b)?)
        end

        consume iso res