use messages = "../messages"
use "../math"

class val Scene
    let objects: Array[Object] val
    let camera: Camera

    new val create(objects': Array[Object] val, camera': Camera) =>
        objects = objects'
        camera = camera'

class val Object
    let shape: Shape
    let material: Material

    new val create(shape': Shape val, material': Material val) =>
        shape = shape'
        material = material'

class val Camera
    let origin: Point3
    let lower_left_corner: Point3
    let horizontal: Vec3
    let vertical: Vec3

    new val create(look_from: Point3, look_at: Point3, vup: Vec3, fov: F64, aspect_ratio: F64) =>
        let theta = Degrees.to_radians(fov)
        let h = (theta/2).tan()
        let viewport_height: F64 = 2.0 * h
        let viewport_width: F64 = aspect_ratio * viewport_height

        let w = (look_from - look_at).unit()
        let u = vup.cross(w).unit()
        let v = w.cross(u)

        origin = look_from
        horizontal = u * viewport_width
        vertical = v * viewport_height
        lower_left_corner = origin - (horizontal/2) - (vertical/2) - w

primitive Transform
    // At the moment this only transforms a scene from the types
    // defined in messages to the internal types. In the future
    // this might do more significant transforms (BVH?)
    fun apply(scene': messages.Scene, config': messages.Config): Scene =>
        let objects' = recover Array[Object](scene'.objects.size()) end

        for o in scene'.objects.values() do
            objects'.push(transform_object(o))
        end

        Scene(consume objects', transform_camera(scene'.camera, config'))

    fun transform_camera(camera': messages.Camera, config': messages.Config): Camera =>
        Camera(where
            look_from = transform_position(camera'.look_from),
            look_at = transform_position(camera'.look_at),
            fov = camera'.fov,
            vup = Vec3(0, 1, 0),
            aspect_ratio = config'.aspect_ratio
        )

    fun transform_object(object': messages.Object): Object =>
        Object(where
            shape' = transform_shape(object'.shape),
            material' = transform_material(object'.material)
        )

    fun transform_shape(shape: messages.Shape): Shape =>
        match shape.one_of
        | let s: messages.Sphere =>
            Sphere(where
                origin' = transform_position(s.origin),
                radius' = s.radius
            )
        end

    fun transform_position(p': messages.Position): Vec3 =>
        Vec3(where
            x' = p'.x,
            y' = p'.y,
            z' = p'.z
        )

    fun transform_material(material: messages.Material): Material =>
        match material.one_of
        | let l: messages.Lambert =>
            Lambert(where
                albedo' = transform_color(l.albedo)
            )
        | let m: messages.Metal =>
            Metal(where
                albedo' = transform_color(m.albedo),
                scatter' = m.scatter
            )
        | let d: messages.Dielectric =>
            Dielectric(where
                ior = d.index_of_refraction
            )
        end

    fun transform_color(color: messages.Color): Vec3 =>
        Vec3(where
            x' = color.r.f64()/255,
            y' = color.g.f64()/255,
            z' = color.b.f64()/255
        )