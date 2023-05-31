// Code generated by submsg; DO NOT EDIT. 
use "msgpack" 
use "buffered" 
use "debug" 
use "collections" 
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
    fun config(): MsgId => 1
    fun render_frame(): MsgId => 2

interface tag CoreServer
    fun tag info(body: Array[U8] iso)
    fun tag config(body: Array[U8] iso)
    fun tag render_frame(body: Array[U8] iso)

primitive CoreRouter
    fun apply(s: CoreServer): ReceiveMsg =>
        {(id: MsgId, body: Array[U8] iso) =>
            match id
            | Core.info() => s.info(consume body)
            | Core.config() => s.config(consume body)
            | Core.render_frame() => s.render_frame(consume body)
            end
        }

actor CoreClient
    let send_msg: SendMsg

    new create(send_msg': SendMsg) =>
        send_msg = send_msg'

    be info(data: (Array[U8 val] iso | None)) =>
        send_msg(Core.info(), consume data)

    be config(data: (Array[U8 val] iso | None)) =>
        send_msg(Core.config(), consume data)

    be render_frame(data: (Array[U8 val] iso | None)) =>
        send_msg(Core.render_frame(), consume data)


class val MsgCoreInfo is MsgPackMarshalable
    var version: String

    new val create(
        version': String
        ) =>
        version = version'

    new val zero() =>
        version = ""

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 1)?
        MessagePackEncoder.fixstr(w, "Version")?
        MessagePackEncoder.str_8(w, version)?

primitive UnmarshalMsgPackMsgCoreInfo
    fun apply(r: Reader ref): MsgCoreInfo =>
        var version': String = ""

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Version" =>
                    version' = MessagePackDecoder.str(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        MsgCoreInfo(
        version'
        )
class val Config is MsgPackMarshalable
    var aspect_ratio: F64
    var depth: U64
    var frame_count: U64
    var frame_speed: U64
    var image_height: U64
    var image_width: U64
    var samples: U64

    new val create(
        aspect_ratio': F64,
        depth': U64,
        frame_count': U64,
        frame_speed': U64,
        image_height': U64,
        image_width': U64,
        samples': U64
        ) =>
        aspect_ratio = aspect_ratio'
        depth = depth'
        frame_count = frame_count'
        frame_speed = frame_speed'
        image_height = image_height'
        image_width = image_width'
        samples = samples'

    new val zero() =>
        aspect_ratio = 0.0
        depth = 0
        frame_count = 0
        frame_speed = 0
        image_height = 0
        image_width = 0
        samples = 0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 7)?
        MessagePackEncoder.fixstr(w, "AspectRatio")?
        MessagePackEncoder.float_64(w, aspect_ratio)
        MessagePackEncoder.fixstr(w, "Depth")?
        MessagePackEncoder.uint_64(w, depth)
        MessagePackEncoder.fixstr(w, "FrameCount")?
        MessagePackEncoder.uint_64(w, frame_count)
        MessagePackEncoder.fixstr(w, "FrameSpeed")?
        MessagePackEncoder.uint_64(w, frame_speed)
        MessagePackEncoder.fixstr(w, "ImageHeight")?
        MessagePackEncoder.uint_64(w, image_height)
        MessagePackEncoder.fixstr(w, "ImageWidth")?
        MessagePackEncoder.uint_64(w, image_width)
        MessagePackEncoder.fixstr(w, "Samples")?
        MessagePackEncoder.uint_64(w, samples)

primitive UnmarshalMsgPackConfig
    fun apply(r: Reader ref): Config =>
        var aspect_ratio': F64 = 0.0
        var depth': U64 = 0
        var frame_count': U64 = 0
        var frame_speed': U64 = 0
        var image_height': U64 = 0
        var image_width': U64 = 0
        var samples': U64 = 0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "AspectRatio" =>
                    aspect_ratio' = MessagePackDecoder.f64(r)?
                | "Depth" =>
                    depth' = MessagePackDecoder.u64(r)?
                | "FrameCount" =>
                    frame_count' = MessagePackDecoder.u64(r)?
                | "FrameSpeed" =>
                    frame_speed' = MessagePackDecoder.u64(r)?
                | "ImageHeight" =>
                    image_height' = MessagePackDecoder.u64(r)?
                | "ImageWidth" =>
                    image_width' = MessagePackDecoder.u64(r)?
                | "Samples" =>
                    samples' = MessagePackDecoder.u64(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Config(
        aspect_ratio',
        depth',
        frame_count',
        frame_speed',
        image_height',
        image_width',
        samples'
        )
class val Position is MsgPackMarshalable
    var x: F64
    var y: F64
    var z: F64

    new val create(
        x': F64,
        y': F64,
        z': F64
        ) =>
        x = x'
        y = y'
        z = z'

    new val zero() =>
        x = 0.0
        y = 0.0
        z = 0.0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 3)?
        MessagePackEncoder.fixstr(w, "X")?
        MessagePackEncoder.float_64(w, x)
        MessagePackEncoder.fixstr(w, "Y")?
        MessagePackEncoder.float_64(w, y)
        MessagePackEncoder.fixstr(w, "Z")?
        MessagePackEncoder.float_64(w, z)

primitive UnmarshalMsgPackPosition
    fun apply(r: Reader ref): Position =>
        var x': F64 = 0.0
        var y': F64 = 0.0
        var z': F64 = 0.0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "X" =>
                    x' = MessagePackDecoder.f64(r)?
                | "Y" =>
                    y' = MessagePackDecoder.f64(r)?
                | "Z" =>
                    z' = MessagePackDecoder.f64(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Position(
        x',
        y',
        z'
        )
class val Color is MsgPackMarshalable
    var b: U8
    var g: U8
    var r: U8

    new val create(
        b': U8,
        g': U8,
        r': U8
        ) =>
        b = b'
        g = g'
        r = r'

    new val zero() =>
        b = 0
        g = 0
        r = 0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 3)?
        MessagePackEncoder.fixstr(w, "B")?
        MessagePackEncoder.uint_8(w, b)
        MessagePackEncoder.fixstr(w, "G")?
        MessagePackEncoder.uint_8(w, g)
        MessagePackEncoder.fixstr(w, "R")?
        MessagePackEncoder.uint_8(w, r)

primitive UnmarshalMsgPackColor
    fun apply(r: Reader ref): Color =>
        var b': U8 = 0
        var g': U8 = 0
        var r': U8 = 0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "B" =>
                    b' = MessagePackDecoder.u8(r)?
                | "G" =>
                    g' = MessagePackDecoder.u8(r)?
                | "R" =>
                    r' = MessagePackDecoder.u8(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Color(
        b',
        g',
        r'
        )
class val Pixel is MsgPackMarshalable
    var color: Color
    var x: U64
    var y: U64

    new val create(
        color': Color,
        x': U64,
        y': U64
        ) =>
        color = color'
        x = x'
        y = y'

    new val zero() =>
        color = Color.zero()
        x = 0
        y = 0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 3)?
        MessagePackEncoder.fixstr(w, "Color")?
        color.marshal_msgpack(w)?
        MessagePackEncoder.fixstr(w, "X")?
        MessagePackEncoder.uint_64(w, x)
        MessagePackEncoder.fixstr(w, "Y")?
        MessagePackEncoder.uint_64(w, y)

primitive UnmarshalMsgPackPixel
    fun apply(r: Reader ref): Pixel =>
        var color': Color = Color.zero()
        var x': U64 = 0
        var y': U64 = 0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Color" =>
                    color' = UnmarshalMsgPackColor(r)
                | "X" =>
                    x' = MessagePackDecoder.u64(r)?
                | "Y" =>
                    y' = MessagePackDecoder.u64(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Pixel(
        color',
        x',
        y'
        )
class val Lambert is MsgPackMarshalable
    var albedo: Color

    new val create(
        albedo': Color
        ) =>
        albedo = albedo'

    new val zero() =>
        albedo = Color.zero()

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 1)?
        MessagePackEncoder.fixstr(w, "Albedo")?
        albedo.marshal_msgpack(w)?

primitive UnmarshalMsgPackLambert
    fun apply(r: Reader ref): Lambert =>
        var albedo': Color = Color.zero()

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Albedo" =>
                    albedo' = UnmarshalMsgPackColor(r)
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Lambert(
        albedo'
        )
class val Metal is MsgPackMarshalable
    var albedo: Color
    var scatter: F64

    new val create(
        albedo': Color,
        scatter': F64
        ) =>
        albedo = albedo'
        scatter = scatter'

    new val zero() =>
        albedo = Color.zero()
        scatter = 0.0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 2)?
        MessagePackEncoder.fixstr(w, "Albedo")?
        albedo.marshal_msgpack(w)?
        MessagePackEncoder.fixstr(w, "Scatter")?
        MessagePackEncoder.float_64(w, scatter)

primitive UnmarshalMsgPackMetal
    fun apply(r: Reader ref): Metal =>
        var albedo': Color = Color.zero()
        var scatter': F64 = 0.0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Albedo" =>
                    albedo' = UnmarshalMsgPackColor(r)
                | "Scatter" =>
                    scatter' = MessagePackDecoder.f64(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Metal(
        albedo',
        scatter'
        )
class val Dielectric is MsgPackMarshalable
    var index_of_refraction: F64

    new val create(
        index_of_refraction': F64
        ) =>
        index_of_refraction = index_of_refraction'

    new val zero() =>
        index_of_refraction = 0.0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 1)?
        MessagePackEncoder.fixstr(w, "IndexOfRefraction")?
        MessagePackEncoder.float_64(w, index_of_refraction)

primitive UnmarshalMsgPackDielectric
    fun apply(r: Reader ref): Dielectric =>
        var index_of_refraction': F64 = 0.0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "IndexOfRefraction" =>
                    index_of_refraction' = MessagePackDecoder.f64(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Dielectric(
        index_of_refraction'
        )
class val Material is MsgPackMarshalable
    var one_of: (
        Lambert |
        Metal |
        Dielectric 
    )

    new val create(
        one_of': (
            Lambert |
            Metal |
            Dielectric 
            )
        ) =>

        one_of = one_of'

    new val zero() =>
        one_of = Lambert.zero()

    fun marshal_msgpack(w: Writer ref)? =>
        match one_of
        | let o: Lambert =>
            MessagePackEncoder.uint_8(w, 0)
            o.marshal_msgpack(w)?
        | let o: Metal =>
            MessagePackEncoder.uint_8(w, 1)
            o.marshal_msgpack(w)?
        | let o: Dielectric =>
            MessagePackEncoder.uint_8(w, 2)
            o.marshal_msgpack(w)?
        end

primitive UnmarshalMsgPackMaterial
    fun apply(r: Reader ref): Material =>
        try

        Material(match MessagePackDecoder.u8(r)?
        | 0 => UnmarshalMsgPackLambert(r)
        | 1 => UnmarshalMsgPackMetal(r)
        | 2 => UnmarshalMsgPackDielectric(r)
        else
            Debug("broken oneof" where stream = DebugErr)
            Lambert.zero()
        end)

        else
            Debug("broken oneof 2" where stream = DebugErr)
            Material.zero()
        end
class val Sphere is MsgPackMarshalable
    var origin: Position
    var radius: F64

    new val create(
        origin': Position,
        radius': F64
        ) =>
        origin = origin'
        radius = radius'

    new val zero() =>
        origin = Position.zero()
        radius = 0.0

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 2)?
        MessagePackEncoder.fixstr(w, "Origin")?
        origin.marshal_msgpack(w)?
        MessagePackEncoder.fixstr(w, "Radius")?
        MessagePackEncoder.float_64(w, radius)

primitive UnmarshalMsgPackSphere
    fun apply(r: Reader ref): Sphere =>
        var origin': Position = Position.zero()
        var radius': F64 = 0.0

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Origin" =>
                    origin' = UnmarshalMsgPackPosition(r)
                | "Radius" =>
                    radius' = MessagePackDecoder.f64(r)?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Sphere(
        origin',
        radius'
        )
class val Shape is MsgPackMarshalable
    var one_of: (
        Sphere 
    )

    new val create(
        one_of': (
            Sphere 
            )
        ) =>

        one_of = one_of'

    new val zero() =>
        one_of = Sphere.zero()

    fun marshal_msgpack(w: Writer ref)? =>
        match one_of
        | let o: Sphere =>
            MessagePackEncoder.uint_8(w, 0)
            o.marshal_msgpack(w)?
        end

primitive UnmarshalMsgPackShape
    fun apply(r: Reader ref): Shape =>
        try

        Shape(match MessagePackDecoder.u8(r)?
        | 0 => UnmarshalMsgPackSphere(r)
        else
            Debug("broken oneof" where stream = DebugErr)
            Sphere.zero()
        end)

        else
            Debug("broken oneof 2" where stream = DebugErr)
            Shape.zero()
        end
class val Object is MsgPackMarshalable
    var material: Material
    var shape: Shape

    new val create(
        material': Material,
        shape': Shape
        ) =>
        material = material'
        shape = shape'

    new val zero() =>
        material = Material.zero()
        shape = Shape.zero()

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 2)?
        MessagePackEncoder.fixstr(w, "Material")?
        material.marshal_msgpack(w)?
        MessagePackEncoder.fixstr(w, "Shape")?
        shape.marshal_msgpack(w)?

primitive UnmarshalMsgPackObject
    fun apply(r: Reader ref): Object =>
        var material': Material = Material.zero()
        var shape': Shape = Shape.zero()

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Material" =>
                    material' = UnmarshalMsgPackMaterial(r)
                | "Shape" =>
                    shape' = UnmarshalMsgPackShape(r)
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Object(
        material',
        shape'
        )
class val Camera is MsgPackMarshalable
    var aperture: F64
    var fov: F64
    var look_at: Position
    var look_from: Position

    new val create(
        aperture': F64,
        fov': F64,
        look_at': Position,
        look_from': Position
        ) =>
        aperture = aperture'
        fov = fov'
        look_at = look_at'
        look_from = look_from'

    new val zero() =>
        aperture = 0.0
        fov = 0.0
        look_at = Position.zero()
        look_from = Position.zero()

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 4)?
        MessagePackEncoder.fixstr(w, "Aperture")?
        MessagePackEncoder.float_64(w, aperture)
        MessagePackEncoder.fixstr(w, "Fov")?
        MessagePackEncoder.float_64(w, fov)
        MessagePackEncoder.fixstr(w, "LookAt")?
        look_at.marshal_msgpack(w)?
        MessagePackEncoder.fixstr(w, "LookFrom")?
        look_from.marshal_msgpack(w)?

primitive UnmarshalMsgPackCamera
    fun apply(r: Reader ref): Camera =>
        var aperture': F64 = 0.0
        var fov': F64 = 0.0
        var look_at': Position = Position.zero()
        var look_from': Position = Position.zero()

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Aperture" =>
                    aperture' = MessagePackDecoder.f64(r)?
                | "Fov" =>
                    fov' = MessagePackDecoder.f64(r)?
                | "LookAt" =>
                    look_at' = UnmarshalMsgPackPosition(r)
                | "LookFrom" =>
                    look_from' = UnmarshalMsgPackPosition(r)
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Camera(
        aperture',
        fov',
        look_at',
        look_from'
        )
class val Scene is MsgPackMarshalable
    var camera: Camera
    var objects: Array[Object] val

    new val create(
        camera': Camera,
        objects': Array[Object] val
        ) =>
        camera = camera'
        objects = objects'

    new val zero() =>
        camera = Camera.zero()
        objects = Array[Object]

    fun marshal_msgpack(w: Writer ref)? =>
        MessagePackEncoder.fixmap(w, 2)?
        MessagePackEncoder.fixstr(w, "Camera")?
        camera.marshal_msgpack(w)?
        MessagePackEncoder.fixstr(w, "Objects")?
        Marshal.array_header(w, objects.size())?
        for item' in objects.values() do
            item'.marshal_msgpack(w)?
        end

primitive UnmarshalMsgPackScene
    fun apply(r: Reader ref): Scene =>
        var camera': Camera = Camera.zero()
        var objects': Array[Object] val = Array[Object]

        try
            let map_size = Unmarshal.map(r)?
            for i in Range(0, map_size) do
                let field_name = MessagePackDecoder.fixstr(r)?
                match field_name
                | "Camera" =>
                    camera' = UnmarshalMsgPackCamera(r)
                | "Objects" =>
                    objects' = Unmarshal.array[Object](r, UnmarshalMsgPackObject~apply())?
                else
                    var error_message = String()
                    error_message.append("unknown field: ")
                    error_message.append(consume field_name)

                    Debug(error_message where stream = DebugErr)
                end
            end
        else
            Debug("Error unmarshalling" where stream = DebugErr)
        end

        Scene(
        camera',
        consume objects'
        )