use fj = "fork_join"
use "runtime_info"

use "../messages"

type PixelLoc is (U64, U64)
type PixelColor is (U8, U8, U8)
type OnPixelComplete is {(U64, U64, PixelColor)} val
type OnRenderComplete is {()} val

type JobInput is PixelLoc
type JobOutput is (PixelLoc, PixelColor)

primitive Renderer
    fun render(auth: SchedulerInfoAuth val, config: Config, on_pixel: OnPixelComplete, on_complete: OnRenderComplete) =>
        fj.Job[JobInput, JobOutput](
            WorkerBuilder(config),
            PixelGenerator(config.image_width, config.image_height),
            RenderTarget(config.image_height, on_pixel, on_complete),
            auth
        ).start()

class WorkerBuilder is fj.WorkerBuilder[JobInput, JobOutput]
    let _config: Config

    new iso create(config: Config) =>
        _config = config

    fun ref apply(): fj.Worker[JobInput, JobOutput] iso^ =>
        RenderWorker(_config)

class PixelGenerator is fj.Generator[JobInput]
    let _iter: PixelIter

    new iso create(width: U64, height: U64) =>
        _iter = PixelIter(width, height)

    fun ref init(workers: USize) =>
        None

    fun ref apply(): JobInput ? =>
        if _iter.has_next() then
            _iter.next()
        else
            error
        end

class RenderTarget is fj.Collector[JobInput, JobOutput]
    let _on_pixel: OnPixelComplete
    let _on_complete: OnRenderComplete
    let _height: U64

    new iso create(height: U64, on_pixel: OnPixelComplete, on_complete: OnRenderComplete) =>
        _on_pixel = on_pixel
        _on_complete = on_complete
        _height = height

    fun ref collect(runner: fj.CollectorRunner[JobInput, JobOutput] ref, result: JobOutput) =>
        (let pixel, let color) = result

        // Image renders upside-down, so flip it
        let x = pixel._1
        let y = (_height-1 - pixel._2)
        _on_pixel(x, y, color)

    fun ref finish() =>
        _on_complete()

class RenderWorker is fj.Worker[JobInput, JobOutput]
    var _tracer: Tracer
    var _pixel: PixelLoc = (0, 0)

    new iso create(config: Config) =>
        _tracer = Tracer(config)

    fun ref receive(pixel: JobInput) =>
        _pixel = pixel

    fun ref process(runner: fj.WorkerRunner[JobInput, JobOutput] ref) =>
        runner.deliver((_pixel, _tracer(_pixel)))