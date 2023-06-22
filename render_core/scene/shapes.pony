use "../math"

interface ShapeVisitor[T: Any]
    fun visit_sphere(s: Sphere box): T
    fun visit_shape_list(s: ShapeList box): T
    fun visit_shape_bvh(s: ShapeBVH box): T

interface val Shape
    fun accept[T: Any](s: ShapeVisitor[T] box): T

class val ShapeList is Shape
    let shapes: Array[Shape] val

    new val create(shapes': Array[Shape] val) =>
        shapes = shapes'

    fun accept[T: Any](s: ShapeVisitor[T] box): T =>
        s.visit_shape_list(this)

class val ShapeBVH is Shape
    let left: Shape
    let right: Shape
    let bounding_box: AABB

    new val create(left': Shape, right': Shape, bounding_box': AABB) =>
        left = left'
        right = right'
        bounding_box = bounding_box'

    fun accept[T: Any](s: ShapeVisitor[T] box): T =>
        s.visit_shape_bvh(this)

class val Sphere is Shape
    let origin: Vec3
    let radius: F64
    let material: Material

    new val create(origin': Vec3, radius': F64, material': Material) =>
        origin = origin'
        radius = radius'
        material = material'

    fun accept[T: Any](s: ShapeVisitor[T] box): T =>
        s.visit_sphere(this)