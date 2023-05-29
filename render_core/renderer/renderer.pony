use fj = "fork_join"
use "runtime_info"

use messages = "../messages"
use "../scene"

// Renderer is the interface between tracer and the outside world
// It orchestrates the render job for one image, calling tracer for every pixel, tracking progress, etc
// This is all plumbing, none of the actual tracing is done here

type PixelLoc is (U64, U64)
type PixelColor is (U8, U8, U8)
type OnPixelComplete is {(U64, U64, PixelColor)} val
type OnRenderComplete is {()} val

type JobInput is PixelLoc
type JobOutput is (PixelLoc, PixelColor)

primitive Renderer
    fun render(auth: SchedulerInfoAuth val, config: messages.Config, scene: Scene, on_pixel: OnPixelComplete, on_complete: OnRenderComplete) =>
        fj.Job[JobInput, JobOutput](
            WorkerBuilder(config, scene),
            PixelGenerator(config.image_width, config.image_height),
            RenderTarget(on_pixel, on_complete),
            auth
        ).start()

class WorkerBuilder is fj.WorkerBuilder[JobInput, JobOutput]
    let _config: messages.Config
    let _scene: Scene

    new iso create(config: messages.Config, scene: Scene) =>
        _config = config
        _scene = scene

    fun ref apply(): fj.Worker[JobInput, JobOutput] iso^ =>
        RenderWorker(_config, _scene)

class PixelGenerator is fj.Generator[JobInput]
    let _width: U64
    let _height: U64

    var _x: U64 = 0
    var _y: U64 = 0

    new iso create(width: U64, height: U64) =>
        _width = width
        _height = height

    fun ref init(workers: USize) =>
        None

    fun ref apply(): JobInput ? =>
        let res: JobInput = (_x, _y)

        if (_x+1) == _width then
            _x = 0
            _y = _y + 1
        else
            _x = _x + 1
        end

        if _y == _height then
            error
        end

        res

class RenderTarget is fj.Collector[JobInput, JobOutput]
    let _on_pixel: OnPixelComplete
    let _on_complete: OnRenderComplete

    new iso create(on_pixel: OnPixelComplete, on_complete: OnRenderComplete) =>
        _on_pixel = on_pixel
        _on_complete = on_complete

    fun ref collect(runner: fj.CollectorRunner[JobInput, JobOutput] ref, result: JobOutput) =>
        (let pixel, let color) = result
        _on_pixel(pixel._1, pixel._2, color)

    fun ref finish() =>
        _on_complete()

class RenderWorker is fj.Worker[JobInput, JobOutput]
    var _tracer: Tracer
    var _pixel: PixelLoc = (0, 0)

    new iso create(config: messages.Config, scene: Scene) =>
        _tracer = Tracer(config, scene)

    fun ref receive(pixel: JobInput) =>
        _pixel = pixel

    fun ref process(runner: fj.WorkerRunner[JobInput, JobOutput] ref) =>
        runner.deliver((_pixel, _tracer(_pixel)))