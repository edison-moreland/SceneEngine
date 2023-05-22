actor Logger
    let out: OutStream tag

    new create(out': OutStream tag) =>
        out = out'

    be log(msg: String) =>
        out.print(msg)