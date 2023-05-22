use fj = "fork_join"
use "files"
use "random"
use "runtime_info"
use "format"

use "../logger"
use "../scene"
use "../math"
use "../messages"

type PixelLoc is (U64, U64)
type OnPixelComplete is {(U64, U64, Color)} val
type OnRenderComplete is {()} val

primitive Renderer
    fun render(auth: SchedulerInfoAuth val, logger: Logger, config: Config, scene: Scene, on_pixel: OnPixelComplete, on_complete: OnRenderComplete) =>
        let job = fj.Job[PixelLoc, (PixelLoc,Vec3)](
            WorkerBuilder(config, scene),
            PixelGenerator(config.image_width, config.image_height),
            RenderTarget(on_pixel, on_complete),
            auth)
        job.start()

class WorkerBuilder is fj.WorkerBuilder[PixelLoc, (PixelLoc,Vec3)]
    let _config: Config
    let _scene: Scene

    new iso create(config: Config, scene: Scene) =>
        _config = config
        _scene = scene

    fun ref apply(): fj.Worker[PixelLoc, (PixelLoc,Vec3)] iso^ =>
        RenderWorker(_config, _scene)

class PixelGenerator is fj.Generator[PixelLoc]
    let _iter: PixelIter

    new iso create(width: U64, height: U64) =>
        _iter = PixelIter(width, height)

    fun ref init(workers: USize) =>
        None

    fun ref apply(): PixelLoc ? =>
        if _iter.has_next() then
            _iter.next()
        else
            error
        end

class RenderTarget is fj.Collector[PixelLoc, (PixelLoc,Vec3)]
    let _on_pixel: OnPixelComplete
    let _on_complete: OnRenderComplete

    new iso create(on_pixel: OnPixelComplete, on_complete: OnRenderComplete) =>
        _on_pixel = on_pixel
        _on_complete = on_complete

    fun clamp(x: F64, min: F64, max: F64): F64 =>
        x.max(min).min(max) //TODO: Is this right?

    fun ref collect(runner: fj.CollectorRunner[PixelLoc, (PixelLoc,Vec3)] ref, result: (PixelLoc,Vec3)) =>
        (let pixel, let color) = result
    
        let color' = Color(
            clamp(color.x.sqrt(), 0.0, 0.999),
            clamp(color.y.sqrt(), 0.0, 0.999),
            clamp(color.z.sqrt(), 0.0, 0.999)
        )

        _on_pixel(pixel._1, pixel._2, color')


    fun ref finish() =>
        _on_complete()

class RenderWorker is fj.Worker[PixelLoc, (PixelLoc,Vec3)]
    let _rand: Rand = _rand.create()

    let _config: Config
    let _scene: Scene

    var _pixel: PixelLoc = (0, 0)

    new iso create(config: Config, scene: Scene) =>
        _config = config
        _scene = scene    

    fun ref receive(pixel: PixelLoc) =>
        _pixel = pixel

    fun ref process(runner: fj.WorkerRunner[PixelLoc, (PixelLoc,Vec3)] ref) =>
        runner.deliver((_pixel, render_pixel(_pixel._1, _pixel._2)))

    fun ref render_pixel(x: U64, y: U64): Vec3 =>
        var color = Vec3.zero()
        var sample: U64 = 0
        while sample < _config.samples do
            let u: F64 = (x.f64() + _rand.real()) / (_config.image_width - 1).f64()
            let v: F64 = (y.f64() + _rand.real()) / (_config.image_height - 1).f64()

            let ray = _scene.camera.get_ray(u, v)

            color = color + trace(ray, _config.depth)

            sample = sample + 1
        end

        color/_config.samples.f64()

    fun ref trace(r: Ray, depth: U64): Vec3 =>
        if depth <= 0 then
            return Vec3.zero()
        end

        match _scene.root_object.hit(r, 0.001, F64.max_value())
        | let rec: HitRecord =>  
            // Bounce ray again
            let next_target: Point3 = rec.p + rec.normal + RandomVec3.unit(_rand)
            trace(Ray(rec.p, next_target - rec.p), depth-1) * 0.5 
        | None => _scene.sky_color(r)
        end