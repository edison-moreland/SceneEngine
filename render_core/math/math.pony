use "random"
use "format"

primitive RandomVec3
    fun range(rand: Rand, min: F64, max: F64): Vec3 =>
        Vec3(
            min + (rand.real() * (max - min)),
            min + (rand.real() * (max - min)),
            min + (rand.real() * (max - min))
        )
    
    fun unit_circle(rand: Rand): Vec3 =>
        // TODO: Redo this
        while true do
            let p = range(rand, -1, 1)
            if (p.length_squared() >= 1) then
                continue
            end
            return p
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

class val Vec3
    let x: F64
    let y: F64
    let z: F64

    new val create(x': F64, y': F64, z': F64) => x = x'; y = y'; z = z'
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

type Point3 is Vec3

class val Ray
    let origin: Point3
    let direction: Vec3

    new val create(origin': Point3, direction': Vec3) =>
        origin = origin'
        direction = direction'

    fun at(t: F64): Point3 =>
        origin + (direction * t) 