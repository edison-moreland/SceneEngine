use "../math"

interface ShapeVisitor[T: Any]
    fun visit_sphere(s: Sphere box): T

interface val Shape
    fun accept[T: Any](s: ShapeVisitor[T] box): T

class val Sphere is Shape
    let origin: Vec3
    let radius: F64

    new val create(origin': Vec3, radius': F64) =>
        origin = origin'
        radius = radius'

    fun accept[T: Any](s: ShapeVisitor[T] box): T =>
        s.visit_sphere(this)