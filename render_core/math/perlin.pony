use "collections"
use "debug"
use "random"

primitive PerlinNoise
    fun new_source(rand: Rand ref): PerlinSource =>
        let size: USize = 256

        PerlinSource(
            _permutation(rand, size),
            _permutation(rand, size),
            _permutation(rand, size),
            _rand_floats(rand, size),
            size
        )

    fun _permutation(rand: Rand ref, size: USize): Array[USize] iso^ =>
        // Initially filled with sequential values
        var ret = recover Array[USize](size) end
        for i in Range(0, size) do
            ret.push(i)
        end

        // Shuffle values around
        for i in Range(size-1, 0) do
            let target = rand.int[USize](size)

            try
                ret.swap_elements(i, target)?
            else
                Debug("perlin machine br0k3" where stream=DebugErr)
            end
        end

        consume ret

    fun _rand_floats(rand: Rand ref, size: USize): Array[F64] iso^ =>
        var ret = recover Array[F64](size) end
        for _ in Range(0, size) do
            ret.push(rand.real())
        end

        consume ret

// TODO: This is currently unfinished, and might stay that way until rendercore get's re-written
class val PerlinSource
    let permutation_size: USize
    let rand_floats: Array[F64]
    let x_permutation: Array[USize]
    let y_permutation: Array[USize]
    let z_permutation: Array[USize]

    new val create(x: Array[USize] iso, y: Array[USize] iso, z: Array[USize] iso, rand_floats': Array[F64] iso, size: USize) =>
        permutation_size = size
        rand_floats = consume rand_floats'
        x_permutation = consume x
        y_permutation = consume y
        z_permutation = consume z

    fun noise(p: Vec3): F64 =>
        let i = (4 * p.x).isize() and (permutation_size - 1).isize()
        let j = (4 * p.y).isize() and (permutation_size - 1).isize()
        let k = (4 * p.z).isize() and (permutation_size - 1).isize()

        try
            rand_floats((i xor j xor k).usize())?
        else
            Debug("perlin machine br0ke" where stream=DebugErr)
            0.0
        end