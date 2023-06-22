
use "random"

primitive RandomVec3
    fun _rand_range(rand: Rand, min: F64, max: F64): F64 =>
        min + (rand.real() * (max - min))

    fun range(rand: Rand, min: F64, max: F64): Vec3 =>
        Vec3(
            _rand_range(rand, min, max),
            _rand_range(rand, min, max),
            _rand_range(rand, min, max)
        )

    fun unit_circle(rand: Rand): Vec3 =>
        while true do
            let p = range(rand, -1, 1)
            if (p.length_squared() >= 1) then
                continue
            end
            return p
        end
        Vec3.zero()

    fun unit_disk(rand: Rand): Vec3 =>
        while true do
            let p = Vec3(where
                x' = _rand_range(rand, -1, 1),
                y' = _rand_range(rand, -1, 1),
                z' = 0
            )

            if p.length_squared() < 1 then
                return p
            end
        end
        Vec3.zero()

    fun unit(rand: Rand): Vec3 =>
        unit_circle(rand).unit()

    fun in_hemisphere(rand: Rand, normal: Vec3): Vec3 =>
        let in_unit_sphere = unit(rand)
        if in_unit_sphere.dot(normal) > 0.0 then
            in_unit_sphere
        else
            -in_unit_sphere
        end

type Point3 is Vec3
class val Vec3
    let x: F64
    let y: F64
    let z: F64

    new val create(x': F64, y': F64, z': F64) => x = x'; y = y'; z = z'
    new val uniform(n: F64) => x = n; y = n; z = n
    new val zero() => x = 0; y = 0; z = 0
    new val one() => x = 1; y = 1; z = 1

    fun apply(i: USize): F64 =>
        match i
        | 0 => x
        | 1 => y
        | 2 => z
        else
            0
        end

    fun add(other: Vec3 box): Vec3 =>
        Vec3(x + other.x, y + other.y, z + other.z)

    fun sub(other: Vec3 box): Vec3 =>
        Vec3(x - other.x, y - other.y, z - other.z)

    fun mul(other: (Vec3 box | F64)): Vec3 =>
        match other
        | let v: Vec3 box => Vec3(x * v.x, y * v.y, z * v.z)
        | let f: F64 => Vec3(x * f, y * f, z * f)
        end

    fun div(other: F64): Vec3 =>
        this * (1/other)

    fun neg(): Vec3 =>
        Vec3(-x, -y, -z)

    fun dot(other: Vec3 box): F64 =>
        (x * other.x) + (y * other.y) + (z * other.z)

    fun cross(other: Vec3 box): Vec3 =>
        Vec3((y * other.z) - (z * other.y),
             (z * other.x) - (x * other.z),
             (x * other.y) - (y * other.x))

    fun unit(): Vec3 =>
        this / this.length()

    fun length_squared(): F64 =>
        this.dot(this)

    fun length(): F64 =>
        length_squared().sqrt()

    fun near_zero(): Bool =>
        let s: F64 = 1e-8
        (x.abs() < s) and (y.abs() < s) and (z.abs() < s)

    fun reflect(n: Vec3): Vec3 =>
        this - (n * this.dot(n) * 2)

    fun refract(n: Vec3, etai_over_etat: F64): Vec3 =>
        let cos_theta: F64 = (-this).dot(n).min(1.0)
        let r_out_perp: Vec3 = (this + (n*cos_theta)) * etai_over_etat
        let r_out_parallel: Vec3 = n * -(1.0 - r_out_perp.length_squared()).abs().sqrt()

        r_out_perp + r_out_parallel

