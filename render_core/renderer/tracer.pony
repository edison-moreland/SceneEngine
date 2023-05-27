use "random"
use "collections"

use messages = "../messages"
use "../scene"
use "../math"

class val _Camera
    let origin: Point3
    let lower_left_corner: Point3
    let horizontal: Vec3
    let vertical: Vec3

    new val create(aspect_ratio: F64) =>
        // let aspect_ratio = 16.0 / 9.0
        let viewport_height: F64 = 2.0
        let viewport_width: F64 = aspect_ratio * viewport_height
        let focal_length: F64 = 1.0

        origin = Vec3.zero()
        horizontal = Vec3(viewport_width, 0.0, 0.0)
        vertical = Vec3(0.0, viewport_height, 0.0)
        lower_left_corner = origin - (horizontal/2) - (vertical/2) - Vec3(0, 0, focal_length)

    fun ray(u: F64, v: F64): Ray =>
        Ray(origin, (lower_left_corner + (horizontal * u) + (vertical * v)) - origin)

class Tracer
    let _rand: Rand = Rand.create()
    let _config: messages.Config
    let _scene: Scene
    let _camera: _Camera

    new create(config: messages.Config, scene: Scene) =>
        _config = config
        _scene = scene
        _camera = _Camera(_config.aspect_ratio)

    fun clamp(x: F64, min: F64, max: F64): F64 =>
        x.max(min).min(max)

    fun to_pixel_color(c: Vec3): PixelColor =>
        (
            (256 * clamp(c.x, 0.0, 0.999)).u8(),
            (256 * clamp(c.y, 0.0, 0.999)).u8(),
            (256 * clamp(c.z, 0.0, 0.999)).u8()
        )

    fun sky_color(r: Ray): Vec3 =>
        let t: F64 = 0.5 * (r.direction.unit().y + 1.0)
        (Vec3(1.0, 1.0, 1.0) * (1.0 - t)) + (Vec3(0.5, 0.7, 1.0) * t )

    fun ref apply(loc: PixelLoc): PixelColor =>
        let x = loc._1
        let y = loc._2

        var color = Vec3.zero()
        for s in Range[U64](0, _config.samples) do
            let u: F64 = (x.f64() + _rand.real()) / (_config.image_width - 1).f64()
            let v: F64 = (y.f64() + _rand.real()) / (_config.image_height - 1).f64()

            let ray = _camera.ray(u, v)

            color = color + trace(ray, _config.depth)

        end

        to_pixel_color(color/_config.samples.f64())

    fun ref trace(r: Ray, depth: U64): Vec3 =>
        if depth <= 0 then
            return Vec3.zero()
        end

        match hit_scene(r)
        | (let rec: HitRecord, let mat: Material) =>
            // Bounce ray again
            match ScatterMaterial(_rand, r, rec)(mat)
            | (let attenuation: Vec3, let scattered: Ray) =>
                attenuation * trace(scattered, depth-1)
            | None =>
                Vec3.zero()
            end
        | None => sky_color(r)
        end

    fun hit_scene(r: Ray): ((HitRecord, Material) | None) =>
        var hit_record: (HitRecord | None) = None
        var closest_t = F64.max_value()
        var material: Material = Lambert(Vec3.zero())

        for o in _scene.objects.values() do
            match HitShape(r, 0.0001, closest_t)(o.shape)
            | let rec: HitRecord =>
                hit_record = rec
                closest_t = rec.t
                material = o.material
            | None => continue
            end
        end

        match hit_record
        | let rec: HitRecord =>
            (rec, material)
        | None =>
            None
        end

class val HitRecord
    let normal: Vec3
    let front_face: Bool
    let p: Point3
    let t: F64

    new val create(normal': Vec3, p': Point3, t': F64) =>
        normal = normal'
        p = p'
        t = t'
        front_face = false

    new val zero() =>
        normal = Vec3.zero()
        front_face = false
        p = Vec3.zero()
        t = 0

    new val from_ray(r: Ray, outward_normal: Vec3,  t': F64, p': Vec3) =>
        t = t'
        p = p'

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
        HitRecord.from_ray(r, (p - s.origin) / s.radius, root, p)

class val ScatterMaterial is MaterialVisitor[((Vec3, Ray) | None)] // (color, attenuation)
    let _rand: Rand ref

    let ray_in: Ray
    let rec: HitRecord

    new create(rand: Rand, ray_in': Ray, rec': HitRecord) =>
        _rand = rand
        ray_in = ray_in'
        rec = rec'

    fun ref apply(m: Material): ((Vec3, Ray) | None) =>
        m.accept[((Vec3, Ray) | None)](this)

    fun ref visit_lambert(l: Lambert box): ((Vec3, Ray) | None) =>
        var scatter_direction = rec.normal + RandomVec3.unit_circle(_rand)

        if scatter_direction.near_zero() then
            scatter_direction = rec.normal
        end

        (
            l.albedo,
            Ray(rec.p, scatter_direction)
        )

    fun ref visit_metal(m: Metal box): ((Vec3, Ray) | None) =>
        let reflected = ray_in.direction.unit().reflect(rec.normal)

        (
            m.albedo,
            Ray(rec.p, reflected)
        )

    fun reflectance(cosine: F64, ri: F64): F64 =>
        // Shlick!
        let r0 = ((1 - ri) / (1 + ri)).pow(2)
        r0 + ((1-r0) * (1 - cosine).pow(5))

    fun ref visit_dielectric(d: Dielectric box): ((Vec3, Ray) | None) =>
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
            Ray(rec.p, direction)
        )