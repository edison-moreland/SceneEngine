use "msgpack"
use "buffered"
use "debug"
use "itertools"
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
                smallest_array(writer, many.size())?
                for m' in many.values() do
                    m'.marshal_msgpack(writer)?
                end
            end
        else
            Debug("error marshalling!" where stream = DebugErr)
        end

        collect(writer.done())

    fun smallest_array(b: Writer ref, size: USize)? =>
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

class val MsgCoreInfo is MsgPackMarshalable
    let version: String

    new val create(version': String) =>
        version = version'

    fun marshal_msgpack(b: Writer ref)? =>
        MessagePackEncoder.fixmap(b, 1)?
        MessagePackEncoder.fixstr(b, "Version")?
        MessagePackEncoder.str_32(b, version)?

class val Config is MsgPackMarshalable
    let image_width: U64
    let image_height: U64

    new val create(image_width': U64, image_height': U64) =>
        image_width = image_width'
        image_height = image_height'

    new val unmarshal_msgpack(b: Array[U8 val] iso) =>
        let rb = recover ref Reader end
        rb.append(consume b)

        var image_width': U64 = 0
        var image_height': U64 = 0

        try
            let map_size = Unmarshal.map(rb)?
            for i in Range(0, map_size) do
                match MessagePackDecoder.fixstr(rb)?
                | "ImageWidth" =>
                    image_width' = MessagePackDecoder.u64(rb)?
                | "ImageHeight" =>
                    image_height' = MessagePackDecoder.u64(rb)?
                else
                    Debug("unknown field" where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        image_width = image_width'
        image_height = image_height'

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 2)?
        MessagePackEncoder.fixstr(w, "ImageWidth")?
        MessagePackEncoder.uint_64(w, image_width)
        MessagePackEncoder.fixstr(w, "ImageHeight")?
        MessagePackEncoder.uint_64(w, image_height)

class val Color is MsgPackMarshalable
    let r: F64
    let g: F64
    let b: F64

    new val create(r': F64, g': F64, b': F64) =>
        r=r'
        g=g'
        b=b'

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 3)?
        MessagePackEncoder.fixstr(w, "R")?
        MessagePackEncoder.float_64(w, r)
        MessagePackEncoder.fixstr(w, "G")?
        MessagePackEncoder.float_64(w, g)
        MessagePackEncoder.fixstr(w, "B")?
        MessagePackEncoder.float_64(w, b)

class val Pixel is MsgPackMarshalable
    let x: U64
    let y: U64
    let color: Color

    new val create(x': U64, y': U64, color': Color) =>
        x=x'
        y=y'
        color=color'

    fun marshal_msgpack(b: Writer ref)? =>
        MessagePackEncoder.fixmap(b, 3)?
        MessagePackEncoder.fixstr(b, "X")?
        MessagePackEncoder.uint_64(b, x)
        MessagePackEncoder.fixstr(b, "Y")?
        MessagePackEncoder.uint_64(b, y)
        MessagePackEncoder.fixstr(b, "Color")?
        color.marshal_msgpack(b)?
