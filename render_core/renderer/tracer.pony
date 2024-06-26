use "random"
use "collections"

use messages = "../messages"
use "../scene"
use "../math"

class Tracer
    let _rand: Rand = Rand.create()
    let _config: messages.Config
    let _scene: Scene

    new create(config: messages.Config, scene: Scene) =>
        _config = config
        _scene = scene

    fun clamp(x: F64, min: F64, max: F64): F64 =>
        x.max(min).min(max)

    fun to_pixel_color(c: Vec3): PixelColor =>
        (
            (256 * clamp(c.x, 0.0, 0.999)).u8(),
            (256 * clamp(c.y, 0.0, 0.999)).u8(),
            (256 * clamp(c.z, 0.0, 0.999)).u8()
        )

    fun sky_color(r: Ray): Vec3 =>
        if _config.sky then
            let t: F64 = 0.5 * (r.direction.unit().y + 1.0)
            (Vec3(1.0, 1.0, 1.0) * (1.0 - t)) + (Vec3(0.5, 0.7, 1.0) * t )
        else
            Vec3.zero()
        end

    fun ref camera_ray(s: F64, t: F64): Ray =>
        let c = _scene.camera

        let rd = RandomVec3.unit_disk(_rand) * c.lens_radius
        let offset = (c.u * rd.x) + (c.v * rd.y)

        Ray(
            (c.origin + offset),
            (c.lower_left_corner + (c.horizontal*s) + (c.vertical*t)) - c.origin - offset
        )

    fun ref apply(loc: PixelLoc): PixelColor =>
        let x = loc._1
        let y = loc._2

        var color = Vec3.zero()
        for i in Range[U64](0, _config.samples) do
            let u: F64 = (x.f64() + _rand.real()) / (_config.image_width - 1).f64()
            let v: F64 = (y.f64() + _rand.real()) / (_config.image_height - 1).f64()

            let ray = camera_ray(u, v)

            color = color + trace(ray, _config.depth)

        end

        to_pixel_color(color/_config.samples.f64())

    fun ref trace(r: Ray, depth: U64): Vec3 =>
        if depth < 1 then
            return Vec3.zero()
        end

        match hit_scene(r)
        | let rec: HitRecord =>
            // Bounce ray again
            match ScatterMaterial(_rand, r, rec)(rec.material)
            | (let attenuation: Vec3, let emitted: Vec3, let scattered: Ray) =>
                emitted + (attenuation * trace(scattered, depth-1))
            | (let emitted: Vec3) =>
                emitted
            | None =>
                Vec3.zero()
            end
        | None => sky_color(r)
        end

    fun hit_scene(r: Ray): (HitRecord | None) =>
        HitShape(r, 0.0001, F64.max_value())(_scene.root_shape)

class val HitRecord
    let material: Material
    let normal: Vec3
    let front_face: Bool
    let p: Point3
    let t: F64
    // Texture coords
    let u: F64
    let v: F64

    new val create(normal': Vec3, p': Point3, t': F64, u': F64, v': F64, material': Material) =>
        normal = normal'
        p = p'
        t = t'
        u = u'
        v = v'
        front_face = false
        material = material'

    new val zero() =>
        normal = Vec3.zero()
        front_face = false
        p = Vec3.zero()
        t = 0
        u = 0
        v = 0
        material = Diffuse(Uniform(Vec3.zero()))

    new val from_ray(r: Ray, outward_normal: Vec3,  t': F64, p': Vec3, u': F64, v': F64, material': Material) =>
        t = t'
        p = p'
        u = u'
        v = v'
        material = material'

        front_face = r.direction.dot(outward_normal) < 0
        normal = if front_face then
            outward_normal
        else
            -outward_normal
        end

class HitShape is ShapeVisitor[(HitRecord | None)]
    let r: Ray
    let t_min: F64
    let t_max: F64

    new create(r': Ray, t_min': F64, t_max': F64) =>
        r = r'
        t_min = t_min'
        t_max = t_max'

    fun apply(s: Shape): (HitRecord | None) =>
        s.accept[(HitRecord | None)](this)

    fun visit_shape_list(s: ShapeList box): (HitRecord | None) =>
        var hit_record: (HitRecord | None) = None
        var closest_t = t_max

        for shape in s.shapes.values() do
            match HitShape(r, t_min, closest_t)(shape)
            | let rec: HitRecord =>
                hit_record = rec
                closest_t = rec.t
            | None => continue
            end
        end

        hit_record

    fun visit_shape_bvh(s: ShapeBVH box): (HitRecord | None) =>
        if not s.bounding_box.hit(r, t_min, t_max) then
            return None
        end

        let hit_left = HitShape(r, t_min, t_max)(s.left)
        let hit_right = HitShape(r, t_min, match hit_left
            | let rec: HitRecord => rec.t
            | None => t_max
            end
        )(s.right)

        match (hit_left, hit_right)
        | (None, None) => None
        | (let le: HitRecord, None) => le
        | (None, let ri: HitRecord) => ri
        | (let le: HitRecord, let ri: HitRecord) =>
            if ri.t < le.t then ri else le end
        end

    fun sphere_uv(p: Vec3): (F64, F64) =>
        let theta = (-p.y).acos()
        let phi = (-p.z).atan2(p.x) + F64.pi()

        (
            phi / (2 * F64.pi()),
            theta / F64.pi()
        )

    fun visit_sphere(s: Sphere box): (HitRecord | None) =>
        let oc = r.origin - s.origin

        let a = r.direction.length_squared()
        let half_b = oc.dot(r.direction)
        let c = oc.length_squared() - (s.radius * s.radius)

        let discriminant = (half_b * half_b) - (a * c)
        if discriminant < 0 then
            return None
        end

        let disc_sqrt = discriminant.sqrt()

        // Find nearest root within range
        var root = (-half_b - disc_sqrt) / a
        if (root < t_min) or (t_max < root) then
            root = (-half_b + disc_sqrt) / a

            if (root < t_min) or (t_max < root) then
                return None
            end
        end

        let p = r.at(root)
        (let u, let v) = sphere_uv(p)
        HitRecord.from_ray(r, (p - s.origin) / s.radius, root, p, u, v, s.material)

type ScatterMaterialReturn is ((Vec3, Vec3, Ray) | Vec3 | None) // (attenuation, emitted, scattered) | (emitted)
class val ScatterMaterial is MaterialVisitor[ScatterMaterialReturn]
    let _rand: Rand ref

    let ray_in: Ray
    let rec: HitRecord

    new create(rand: Rand, ray_in': Ray, rec': HitRecord) =>
        _rand = rand
        ray_in = ray_in'
        rec = rec'

    fun ref apply(m: Material): ScatterMaterialReturn =>
        m.accept[ScatterMaterialReturn](this)

    fun ref visit_diffuse(d: Diffuse box): ScatterMaterialReturn =>
        var scatter_direction = rec.normal + RandomVec3.unit_circle(_rand)

        if scatter_direction.near_zero() then
            scatter_direction = rec.normal
        end

        (
            QueryTexture(rec.u, rec.v, rec.p)(d.texture),
            Vec3.zero(),
            Ray(rec.p, scatter_direction)
        )

    fun ref visit_metallic(m: Metallic box): ScatterMaterialReturn =>
        let reflected = ray_in.direction.unit().reflect(rec.normal)

        (
            QueryTexture(rec.u, rec.v, rec.p)(m.texture),
            Vec3.zero(),
            Ray(rec.p, reflected)
        )

    fun reflectance(cosine: F64, ri: F64): F64 =>
        // Shlick!
        let r0 = ((1 - ri) / (1 + ri)).pow(2)
        r0 + ((1-r0) * (1 - cosine).pow(5))

    fun ref visit_dielectric(d: Dielectric box): ScatterMaterialReturn =>
        let refraction_ratio = if rec.front_face then
            (1/d.index_of_refraction)
        else
            d.index_of_refraction
        end

        let unit_direction = ray_in.direction.unit()

        let cos_theta = (-unit_direction).dot(rec.normal).min(1.0)
        let sin_theta = (1.0 - cos_theta.pow(2)).sqrt()

        let cannot_refract = (refraction_ratio * sin_theta) > 1.0
        let should_reflect = reflectance(cos_theta, refraction_ratio) > _rand.real()

        let direction = if (cannot_refract or should_reflect) then
            unit_direction.reflect(rec.normal)
        else
            unit_direction.refract(rec.normal, refraction_ratio)
        end

        (
            Vec3.one(),
            Vec3.zero(),
            Ray(rec.p, direction)
        )

    fun ref visit_emissive(e: Emissive box): ScatterMaterialReturn =>
        QueryTexture(rec.u, rec.v, rec.p)(e.texture) * e.brightness

class val QueryTexture is TextureVisitor[Vec3]
    let u: F64
    let v: F64
    let p: Vec3

    new create(u': F64, v': F64, p': Vec3) =>
        u = u'
        v = v'
        p = p'

    fun ref apply(t: Texture): Vec3 =>
        t.accept[Vec3](this)

    fun ref visit_uniform(n: Uniform box): Vec3 =>
        n.color

    fun ref visit_checker(c: Checker box): Vec3 =>
        let sines = (10*p.x).sin() * (10*p.y).sin() * (10*p.z).sin()
        if sines < 0 then
            this(c.odd)
        else
            this(c.even)
        end

    fun ref visit_perlin(r: Perlin box): Vec3 =>
        Vec3.one() * r.source.noise(p)