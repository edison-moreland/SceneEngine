use "../math"

interface MaterialVisitor[T: Any]
    fun ref visit_diffuse(d: Diffuse box): T
    fun ref visit_metallic(m: Metallic box): T
    fun ref visit_dielectric(d: Dielectric box): T
    fun ref visit_emissive(e: Emissive box): T

interface val Material
    fun accept[T: Any](m: MaterialVisitor[T] ref): T

class val Diffuse is Material
    let texture: Texture

    new val create(texture': Texture) =>
        texture = texture'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_diffuse(this)

class val Metallic is Material
    let texture: Texture
    let scatter: F64

    new val create(texture': Texture, scatter': F64) =>
        texture = texture'
        scatter = scatter'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_metallic(this)

class val Dielectric is Material
    let index_of_refraction: F64

    new val create(ior: F64) =>
        index_of_refraction = ior

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_dielectric(this)

class val Emissive is Material
    let texture: Texture
    let brightness: F64

    new val create(texture': Texture, brightness': F64) =>
        texture = texture'
        brightness = brightness'

    fun accept[T: Any](m: MaterialVisitor[T] ref): T =>
        m.visit_emissive(this)
