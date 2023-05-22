use "../math"

class val Camera
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

    fun get_ray(u: F64, v: F64): Ray =>
        Ray(origin, (lower_left_corner + (horizontal * u) + (vertical * v)) - origin)