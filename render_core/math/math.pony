use "random"
use "collections"

primitive Degrees
    fun to_radians(d: F64): F64 =>
        (d * F64.pi()) / 180

class val Ray
    let origin: Point3
    let direction: Vec3

    new val create(origin': Point3, direction': Vec3) =>
        origin = origin'
        direction = direction'

    fun at(t: F64): Point3 =>
        origin + (direction * t)

class val AABB // AxisAlignedBoundingBox
    let minimum: Point3
    let maximum: Point3

    new val create(a: Point3, b: Point3) =>
        minimum = a
        maximum = b

    new val zero() =>
        minimum = Point3.zero()
        maximum = Point3.zero()

    new val surrounding(a: AABB, b: AABB) =>
        minimum = Vec3(a.minimum.x.min(b.minimum.x),
                         a.minimum.y.min(b.minimum.y),
                         a.minimum.z.min(b.minimum.z))

        maximum = Vec3(a.maximum.x.max(b.maximum.x),
                       a.maximum.y.max(b.maximum.y),
                       a.maximum.z.max(b.maximum.z))

    fun hit(ray: Ray, t_min': F64, t_max': F64): Bool =>
        // TODO: Try Andrew Kensler's method
        var t_min = t_min'
        var t_max = t_max'

        for i in Range[USize](0, 3) do
            let t0 = ((minimum(i) - ray.origin(i)) / ray.direction(i)).min(
                      (maximum(i) - ray.origin(i)) / ray.direction(i))

            let t1 = ((minimum(i) - ray.origin(i)) / ray.direction(i)).max(
                      (maximum(i) - ray.origin(i)) / ray.direction(i))

            t_min = t0.max(t_min)
            t_max = t1.min(t_max)
            if (t_max <= t_min) then
                return false
            end
        end

        true
