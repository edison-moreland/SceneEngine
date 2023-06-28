use "../math"

interface TextureVisitor[T: Any]
    fun ref visit_uniform(u: Uniform box): T
    fun ref visit_checker(c: Checker box): T
    fun ref visit_perlin(p: Perlin box): T

interface val Texture
    fun accept[T: Any](v: TextureVisitor[T] ref): T

class val Uniform is Texture
    let color: Vec3

    new val create(color': Vec3) =>
        color = color'

    fun accept[T: Any](v: TextureVisitor[T] ref): T =>
        v.visit_uniform(this)

class val Checker is Texture
    let even: Texture
    let odd: Texture

    new val create(even': Texture, odd': Texture) =>
        even = even'
        odd = odd'

    fun accept[T: Any](v: TextureVisitor[T] ref): T =>
        v.visit_checker(this)

class val Perlin is Texture
    let source: PerlinSource
    let scale: F64

    new val create(source': PerlinSource, scale': F64) =>
        source = source'
        scale = scale'

    fun accept[T: Any](v: TextureVisitor[T] ref): T =>
        v.visit_perlin(this)