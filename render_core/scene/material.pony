use "../math"

interface MaterialVisitor
    fun visit_lambert(l: Lambert box)
    fun visit_metal(m: Metal box)
    fun visit_dielectric(d: Dielectric box)

interface val Material
    fun accept(m: MaterialVisitor)

class val Lambert
    let albedo: Vec3

    new val create(albedo': Vec3) =>
        albedo = albedo'

    fun accept(m: MaterialVisitor) =>
        m.visit_lambert(this)

class val Metal
    let albedo: Vec3
    let scatter: F64

    new val create(albedo': Vec3, scatter': F64) =>
        albedo = albedo'
        scatter = scatter'

    fun accept(m: MaterialVisitor) =>
        m.visit_metal(this)

class val Dielectric
    let index_of_refraction: F64

    new val create(ior: F64) =>
        index_of_refraction = ior

    fun accept(m: MaterialVisitor) =>
        m.visit_dielectric(this)
