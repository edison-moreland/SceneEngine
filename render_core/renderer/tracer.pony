use "random"

use messages = "../messages"
use "../scene"

class Tracer
    let _rand: Rand = Rand.create()
    let _config: messages.Config
    let _scene: Scene

    new create(config: messages.Config, scene: Scene) =>
        _config = config
        _scene = scene

    fun clamp(x: F64, min: F64, max: F64): F64 =>
        x.max(min).min(max)

    fun color(r: F64, g: F64, b: F64): PixelColor =>
        (
            (256 * clamp(r, 0.0, 0.999)).u8(),
            (256 * clamp(g, 0.0, 0.999)).u8(),
            (256 * clamp(b, 0.0, 0.999)).u8()
        )

    fun apply(loc: PixelLoc): PixelColor =>
        color(255, 0, 0)
