use "../math"

interface ShapeVisitor
    fun visit_sphere(s: Sphere box)

interface val Shape
    fun accept(s: ShapeVisitor)

class val Sphere is Shape
    let origin: Vec3
    let radius: F64

    new val create(origin': Vec3, radius': F64) =>
        origin = origin'
        radius = radius'

    fun accept(s: ShapeVisitor) =>
        s.visit_sphere(this)