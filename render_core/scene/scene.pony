use messages = "../messages"
use "../math"

class val Scene
    let objects: Array[Object] val

    new val create(objects': Array[Object] val) =>
        objects = objects'

class val Object
    let shape: Shape
    let material: Material

    new val create(shape': Shape val, material': Material val) =>
        shape = shape'
        material = material'

primitive Transform
    // At the moment this only transforms a scene from the types
    // defined in messages to the internal types. In the future
    // this might do more significant transforms (BVH?)
    fun apply(scene': messages.Scene): Scene =>
        let objects' = recover Array[Object](scene'.objects.size()) end

        for o in scene'.objects.values() do
            objects'.push(transformObject(o))
        end

        Scene(consume objects')

    fun transformObject(object': messages.Object): Object =>
        Object(where
            shape' = transformShape(object'.shape),
            material' = transformMaterial(object'.material)
        )

    fun transformShape(shape: messages.Shape): Shape =>
        match shape.one_of
        | let s: messages.Sphere =>
            Sphere(where
                origin' = transformPosition(s.origin),
                radius' = s.radius
            )
        end

    fun transformPosition(p': messages.Position): Vec3 =>
        Vec3(where
            x' = p'.x,
            y' = p'.y,
            z' = p'.z
        )

    fun transformMaterial(material: messages.Material): Material =>
        match material.one_of
        | let l: messages.Lambert =>
            Lambert(where
                albedo' = transformColor(l.albedo)
            )
        | let m: messages.Metal =>
            Metal(where
                albedo' = transformColor(m.albedo),
                scatter' = m.scatter
            )
        | let d: messages.Dielectric =>
            Dielectric(where
                ior = d.index_of_refraction
            )
        end

    fun transformColor(color: messages.Color): Vec3 =>
        Vec3(where
            x' = color.r.f64()/255,
            y' = color.g.f64()/255,
            z' = color.b.f64()/255
        )