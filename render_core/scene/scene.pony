use "../math"

class val SScene
    let root_object: Hittable
    let camera: CCamera

    new val create(root_object': Hittable, camera': CCamera) =>
        root_object = root_object'
        camera = camera'

    fun sky_color(r: Ray): Vec3 =>
        let t: F64 = 0.5 * (r.direction.unit().y + 1.0)
        (Vec3(1.0, 1.0, 1.0) * (1.0 - t)) + (Vec3(0.5, 0.7, 1.0) * t )
