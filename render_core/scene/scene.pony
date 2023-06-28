use "debug"
use "random"

use messages = "../messages"
use "../math"

class val Scene
    let root_shape: Shape
    let camera: Camera

    new val create(root_shape': Shape, camera': Camera) =>
        root_shape = root_shape'
        camera = camera'

class val Camera
    let origin: Point3
    let lower_left_corner: Point3
    let horizontal: Vec3
    let vertical: Vec3
    let u: Vec3
    let v: Vec3
    let w: Vec3
    let lens_radius: F64

    new val create(
        look_from: Point3,
        look_at: Point3,
        vup: Vec3,
        fov: F64,
        aspect_ratio: F64,
        aperture: F64
    ) =>
        let focus_dist = (look_from - look_at).length()

        let theta = Degrees.to_radians(fov)
        let h = (theta/2).tan()
        let viewport_height: F64 = 2.0 * h
        let viewport_width: F64 = aspect_ratio * viewport_height

        w = (look_from - look_at).unit()
        u = vup.cross(w).unit()
        v = w.cross(u)

        origin = look_from
        horizontal = u * viewport_width * focus_dist
        vertical = v * viewport_height * focus_dist
        lower_left_corner = origin - (horizontal/2) - (vertical/2) - (w*focus_dist)

        lens_radius = aperture/2

primitive Transform
    // Transform the scene from the wire into a structure we can trace with
    fun apply(rand: Rand ref, scene': messages.Scene, config': messages.Config): Scene =>
        let shapes' = recover Array[Shape](scene'.objects.size()) end

        for o in scene'.objects.values() do
            // TODO: Objects with multiple shapes
            shapes'.push(transform_object(rand, o))
        end

        let root_shape = if config'.use_bvh then
            try
                ConstructBVH(consume shapes')?
            else
                Debug("BVH construction FUCKED" where stream=DebugErr)
                ShapeList(recover Array[Shape]() end)
            end
        else
            ShapeList(consume shapes')
        end

        Scene(root_shape, transform_camera(scene'.camera, config'))

    fun transform_camera(camera': messages.Camera, config': messages.Config): Camera =>
        Camera(where
            look_from = transform_position(camera'.look_from),
            look_at = transform_position(camera'.look_at),
            fov = camera'.fov,
            aperture = camera'.aperture,
            vup = Vec3(0, 1, 0),
            aspect_ratio = config'.aspect_ratio
        )

    fun transform_object(rand: Rand ref, object': messages.Object): Shape =>
        transform_shape(object'.shape, transform_material(rand, object'.material))

    fun transform_shape(shape: messages.Shape, material: Material): Shape =>
        match shape.one_of
        | let s: messages.Sphere =>
            Sphere(where
                origin' = transform_position(s.origin),
                radius' = s.radius,
                material' = material
            )
        end

    fun transform_position(p': messages.Position): Vec3 =>
        Vec3(where
            x' = p'.x,
            y' = p'.y,
            z' = p'.z
        )

    fun transform_material(rand: Rand ref, material: messages.Material): Material =>
        match material.one_of
        | let l: messages.Diffuse =>
            Diffuse(where
                texture' = transform_texture(rand, l.texture)
            )
        | let m: messages.Metallic =>
            Metallic(where
                texture' = transform_texture(rand, m.texture),
                scatter' = m.scatter
            )
        | let d: messages.Dielectric =>
            Dielectric(where
                ior = d.index_of_refraction
            )
        end

    fun transform_texture(rand: Rand ref, texture: messages.Texture): Texture =>
        match texture.one_of
        | let u: messages.Uniform =>
            Uniform(where
                color' = transform_color(u.color)
            )
        | let c: messages.Checker =>
            Checker(where
                even' = transform_texture(rand, c.even),
                odd' = transform_texture(rand, c.odd)
            )
        | let p: messages.Perlin =>
            Perlin(where
                source' = PerlinNoise.new_source(rand),
                scale' = p.scale
            )
        end

    fun transform_color(color: messages.Color): Vec3 =>
        Vec3(where
            x' = color.r.f64()/255,
            y' = color.g.f64()/255,
            z' = color.b.f64()/255
        )