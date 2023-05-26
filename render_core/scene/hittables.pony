use "../math"

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

    new val from_ray(r: Ray, outward_normal: Vec3,  t': F64, p': Vec3) =>
        t = t'
        p = p'

        front_face = r.direction.dot(outward_normal) < 0
        normal = if front_face then 
            outward_normal 
        else 
            -outward_normal 
        end

trait val Hittable
    fun hit(r: Ray, t_min: F64, t_max: F64): (HitRecord | None)

class val HittableList is Hittable
    let hittables: Array[Hittable] val

    new val create(hittables': Array[Hittable] val) =>
        hittables = hittables'

    fun hit(r: Ray, t_min: F64, t_max: F64): (HitRecord | None) =>
        var hit_record: (HitRecord | None) = None
        var closest_t = t_max

        for hittable in hittables.values() do
            match hittable.hit(r, t_min, closest_t)
            | let rec: HitRecord => 
                hit_record = rec
                closest_t = rec.t
            | None => continue
            end
        end

        hit_record


class val HSphere is Hittable
    let center: Point3
    let radius: F64

    new val create(center': Point3, radius': F64) =>
        center = center'
        radius = radius'
    
    fun hit(r: Ray, t_min: F64, t_max: F64): (HitRecord | None) =>
        let oc = r.origin - center

        let a = r.direction.length_squared()
        let half_b = oc.dot(r.direction)
        let c = oc.length_squared() - (radius * radius)

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
        HitRecord.from_ray(r, (p - center) / radius, root, p)