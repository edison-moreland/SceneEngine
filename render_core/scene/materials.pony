use "../math"

interface MaterialVisitor[T: Any]
    fun ref visit_diffuse(d: Diffuse box): T
    fun ref visit_metallic(m: Metallic box): T
    fun ref visit_dielectric(d: Dielectric box): T

interface val Material
    fun accept[T: Any](m: MaterialVisitor[T] ref): T

class val Diffuse is Material
    let albedo: Vec3

    new val create(albedo': Vec3) =>
        albedo = albedo'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_diffuse(this)

class val Metallic is Material
    let albedo: Vec3
    let scatter: F64

    new val create(albedo': Vec3, scatter': F64) =>
        albedo = albedo'
        scatter = scatter'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_metallic(this)

class val Dielectric is Material
    let index_of_refraction: F64

    new val create(ior: F64) =>
        index_of_refraction = ior

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_dielectric(this)
