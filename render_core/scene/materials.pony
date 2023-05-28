use "../math"

interface MaterialVisitor[T: Any]
    fun ref visit_lambert(l: Lambert box): T
    fun ref visit_metal(m: Metal box): T
    fun ref visit_dielectric(d: Dielectric box): T

interface val Material
    fun accept[T: Any](m: MaterialVisitor[T] ref): T

class val Lambert is Material
    let albedo: Vec3

    new val create(albedo': Vec3) =>
        albedo = albedo'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_lambert(this)

class val Metal is Material
    let albedo: Vec3
    let scatter: F64

    new val create(albedo': Vec3, scatter': F64) =>
        albedo = albedo'
        scatter = scatter'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_metal(this)

class val Dielectric is Material
    let index_of_refraction: F64

    new val create(ior: F64) =>
        index_of_refraction = ior

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_dielectric(this)
