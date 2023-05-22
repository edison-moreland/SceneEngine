class PixelIter
    let width: U64
    let height: U64
    let start_x: U64
    let start_y: U64
    var x: U64
    var y: U64

    new create(width': U64, height': U64) =>
        width = width'
        height = height'
        start_x = 0
        start_y = 0
        x = start_x
        y = start_y

    new chunk(chunk': (U64, U64, U64, U64)) =>
        (let x', let y', let width', let height') = chunk'
        width = width'
        height = height'
        start_x = x'
        start_y = y'
        x = start_x
        y = start_y

    fun has_next(): Bool =>
        (x < (start_x+width)) and (y < (start_y+height))

    fun ref next(): (U64, U64) =>
        let pixel = (x, y)

        x = x + 1
        if x >= (start_x+width) then
            x = start_x
            y = y + 1
        end

        pixel