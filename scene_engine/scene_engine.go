package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"go.uber.org/zap"

	"github.com/edison-moreland/SceneEngine/scene_engine/core"
	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
	"github.com/edison-moreland/SceneEngine/scene_engine/layout"
	"github.com/edison-moreland/SceneEngine/scene_engine/scenebuilder"
	"github.com/edison-moreland/SceneEngine/scene_engine/script"
)

var (
	corePath   string
	scriptPath string
)

var backgroundColor = rl.LightGray

func init() {
	flag.StringVar(&corePath, "core", "", "Path to rendercore")
	flag.StringVar(&scriptPath, "script", "./scene.tengo", "Path to scene script")
	flag.Parse()

	if corePath == "" {
		fmt.Println("parameter `core` is required.")
		os.Exit(1)
	}
}

func main() {
	engineCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("Welcome to SceneEngine!")

	scriptFileWatcher, err := WatchFile(scriptPath)
	if err != nil {
		logger.Fatal("Could not start file watcher!", zap.Error(err))
	}

	logger.Info("Starting render core")
	renderCore, err := core.Start(engineCtx, corePath)
	if err != nil {
		logger.Fatal("Could not start core!", zap.Error(err))
	}
	logger.Info("Render core ready!", zap.String("version", renderCore.Info()))

	rl.InitWindow(400, 400, "SceneEngine")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	var sceneCache scenebuilder.SceneCache
	err = RunPhases(logger, LoadScript, map[AppPhaseId]AppPhase{
		Preview:    PreviewPhase(scriptFileWatcher, &sceneCache),
		LoadScript: LoadScriptPhase(&sceneCache, scriptPath),
		Render:     RenderPhase(logger, renderCore, &sceneCache, exportDir(scriptPath)),
		Encode:     EncodePhase(&sceneCache, exportDir(scriptPath)),
	})
	if err != nil {
		logger.Fatal("Render machine br0k3", zap.Error(err))
	}
}

func exportDir(scriptPath string) string {
	baseDir, scriptName := path.Split(scriptPath)
	name := strings.SplitN(scriptName, ".", 2)[0]
	return path.Join(baseDir, name)
}

type AppPhaseId int

const (
	Preview AppPhaseId = iota
	LoadScript
	Render
	Encode
)

type AppPhase interface {
	Think() (next AppPhaseId, err error)
	Draw()

	Start() error // Called when transitioning to this phase
	End() error   // Called when transitioning to a different phase

	Shutdown() error
}

type emptyPhase struct {
}

func (e *emptyPhase) Start() error {
	return nil
}

func (e *emptyPhase) End() error {
	return nil
}

func (e *emptyPhase) Shutdown() error {
	return nil
}

func RunPhases(logger *zap.Logger, startPhase AppPhaseId, phases map[AppPhaseId]AppPhase) error {
	current := startPhase

	defer func() {
		for _, phase := range phases {
			if err := phase.Shutdown(); err != nil {
				logger.Error("Err completing phase", zap.Error(err))
			}
		}
	}()

	if err := phases[current].Start(); err != nil {
		return err
	}

	for !rl.WindowShouldClose() {
		next, err := phases[current].Think()
		if err != nil {
			return err
		}

		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)
		phases[current].Draw()
		rl.EndDrawing()

		if next != current {
			if err := phases[current].End(); err != nil {
				return err
			}
			if err := phases[next].Start(); err != nil {
				return err
			}
		}
		current = next
	}

	if err := phases[current].End(); err != nil {
		return err
	}

	return nil
}

/* Preview Phase */
// TODO: Maybe rename to "preview"?

type preview struct {
	emptyPhase
	script *FileWatcher

	renderButton   Button
	timelineSlider Slider

	sceneCache   *scenebuilder.SceneCache
	currentScene uint64
}

func PreviewPhase(script *FileWatcher, sceneCache *scenebuilder.SceneCache) AppPhase {
	return &preview{
		script: script,

		renderButton: NewButton("Render", rl.Gray, layout.Window().Margin(5), layout.NorthWest),

		sceneCache:   sceneCache,
		currentScene: 0,
	}
}

func (p *preview) textureColor(t messages.Texture) rl.Color {
	// TODO: use actual textures
	switch t := t.OneOf.(type) {
	case messages.Uniform:
		c := t.Color
		return rl.NewColor(c.R, c.G, c.B, 255)
	case messages.Checker:
		return p.textureColor(t.Even)
	case messages.Perlin:
		return rl.Gray
	}
	panic("Texture not implemented " + fmt.Sprintf("%#v", t))
}

func (p *preview) positionToVec3(pos messages.Position) rl.Vector3 {
	return rl.NewVector3(float32(pos.X), float32(pos.Y), float32(pos.Z))
}

func (p *preview) drawScenePreview() {
	scene := p.sceneCache.Scene(p.currentScene)

	camera := scene.Camera
	rl.BeginMode3D(rl.Camera3D{
		Position: p.positionToVec3(camera.LookFrom),
		Target:   p.positionToVec3(camera.LookAt),
		Up:       rl.NewVector3(0, 1, 0),
		Fovy:     float32(camera.Fov),
	})

	for _, object := range scene.Objects {
		objectColor := rl.Gray
		switch m := object.Material.OneOf.(type) {
		case messages.Diffuse:
			objectColor = p.textureColor(m.Texture)
		case messages.Metallic:
			objectColor = p.textureColor(m.Texture)
		case messages.Emissive:
			objectColor = p.textureColor(m.Texture)

			// Dielectric gets default color
		}

		switch s := object.Shape.OneOf.(type) {
		case messages.Sphere:
			rl.DrawSphere(p.positionToVec3(s.Origin), float32(s.Radius), objectColor)
		}
	}

	rl.EndMode3D()
}

func (p *preview) Start() error {
	config := p.sceneCache.Config()
	rl.SetTargetFPS(int32(config.FrameSpeed))

	sliderRec := layout.
		Window().
		Margin(10).
		Resize(layout.Vertical, layout.South, 40)

	p.timelineSlider = NewSlider(sliderRec, 1, float32(config.FrameCount))

	return nil
}

func (p *preview) Think() (AppPhaseId, error) {
	if p.renderButton.Down() {
		return Render, nil
	}

	if p.script.HasChanged() {
		return LoadScript, nil
	}

	config := p.sceneCache.Config()
	p.currentScene += 1
	if p.currentScene > config.FrameCount {
		p.currentScene = 1
	}
	p.timelineSlider.Update(float32(p.currentScene))

	return Preview, nil
}

func (p *preview) End() error {
	p.script.ClearChange()

	rl.SetTargetFPS(60)
	return nil
}

func (p *preview) Draw() {
	p.drawScenePreview()

	p.renderButton.Draw()
	p.timelineSlider.Draw()
}

/* Render Phase */

type render struct {
	emptyPhase
	logger *zap.Logger

	core       *core.RenderCore
	target     *renderTarget
	sceneCache *scenebuilder.SceneCache
	exportPath string

	renderActive bool

	currentFrame   uint64
	lastFrameStart time.Time
	frameElapsed   rollingAverage
}

func RenderPhase(logger *zap.Logger, core *core.RenderCore, sceneCache *scenebuilder.SceneCache, exportDir string) AppPhase {
	r := render{
		logger: logger,

		core:       core,
		sceneCache: sceneCache,
		exportPath: exportDir,

		renderActive: false,

		currentFrame: 0,
		frameElapsed: rollingAverage{sampleSize: 10},
	}

	return &r
}

func (r *render) prepareExportDir() error {
	if err := os.RemoveAll(r.exportPath); err != nil {
		return err
	}

	return os.Mkdir(r.exportPath, os.FileMode(0777))
}

func (r *render) startFrame() {
	r.logger.Info("Starting render", zap.Uint64("frame", r.currentFrame))
	scene := r.sceneCache.Scene(r.currentFrame)

	r.lastFrameStart = time.Now()
	r.core.WaitForReady() // This should return immediately
	r.target.Reset(r.core.StartRender(scene))
}

func (r *render) Start() error {
	// Begin the render cycle
	r.currentFrame = 0
	r.renderActive = false
	r.frameElapsed.Reset()

	if err := r.prepareExportDir(); err != nil {
		return err
	}

	config := r.sceneCache.Config()
	r.target = newRenderTarget(config.ImageWidth, config.ImageHeight)
	r.core.SetConfig(config)

	return nil
}

func (r *render) Think() (AppPhaseId, error) {
	if !r.renderActive {
		r.currentFrame += 1

		config := r.sceneCache.Config()
		if r.currentFrame > config.FrameCount {
			if config.FrameSpeed > 0 {
				return Encode, nil
			}
			return Preview, nil
		}

		r.startFrame()
		r.renderActive = true
	}

	r.target.RenderBufferToTexture()
	if r.target.done {
		frameTime := time.Since(r.lastFrameStart)
		r.logger.Info("Finished frame", zap.Duration("frame_time", frameTime))

		r.frameElapsed.Sample(frameTime)
		r.target.Export(r.exportPath, r.currentFrame)
		r.renderActive = false
	}

	return Render, nil
}

func (r *render) End() error {
	// There should be a pending ready already
	r.core.WaitForReady()
	r.target.Close()
	r.target = nil
	return nil
}

func (r *render) Draw() {
	rl.DrawTextureV(r.target.Texture, rl.Vector2Zero(), rl.White)

	config := r.sceneCache.Config()
	statBlock := NewStatBlock(layout.Window().Margin(5), layout.NorthWest)
	statBlock.AddStat("frame", fmt.Sprintf("%d/%d", r.currentFrame, config.FrameCount))
	if r.frameElapsed.HasSamples() {
		averageFrameTime := r.frameElapsed.Average()
		statBlock.AddStat("avg_frame_time", averageFrameTime.Round(time.Millisecond).String())
		statBlock.AddStat("done_in", (time.Duration(config.FrameCount-r.currentFrame) * averageFrameTime).Round(time.Millisecond).String())
	}

	statBlock.Draw()
}

func (r *render) Shutdown() error {
	return nil
}

/* Encode Phase */

type encode struct {
	emptyPhase

	sceneCache *scenebuilder.SceneCache

	encodeCmd    *exec.Cmd
	encodeActive bool
	exportDir    string
}

func EncodePhase(sceneCache *scenebuilder.SceneCache, exportDir string) AppPhase {
	return &encode{
		sceneCache: sceneCache,
		exportDir:  exportDir,
	}
}

func (e *encode) startEncode() error {
	inputImages := path.Join(e.exportDir, "frame-%d.png")
	outputVideo := path.Clean(path.Join(e.exportDir, "..", path.Base(e.exportDir)+".mp4"))

	// These ffmpeg parameters should create an MP4 compatible with quicktime
	config := e.sceneCache.Config()
	e.encodeCmd = exec.Command("ffmpeg",
		"-y",
		"-framerate", strconv.Itoa(int(config.FrameSpeed)),
		"-i", inputImages,
		"-pix_fmt", "yuv420p",
		"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
		"-vcodec", "h264",
		"-acodec", "aac",
		outputVideo,
	)
	e.encodeCmd.Stderr = os.Stderr
	e.encodeCmd.Stdout = os.Stdout

	err := e.encodeCmd.Start()
	if err != nil {
		return err
	}

	e.encodeActive = true
	go func() {
		if err := e.encodeCmd.Wait(); err != nil {
			panic(err)
		}
		e.encodeActive = false
	}()

	return nil
}

func (e *encode) Start() error {
	return e.startEncode()
}

func (e *encode) Think() (AppPhaseId, error) {
	if e.encodeActive {
		return Encode, nil
	}

	return Preview, nil
}

func (e *encode) End() error {
	e.encodeCmd = nil

	return nil
}

func (e *encode) Draw() {
	rl.DrawText("Encoding video...", 10, 10, 30, rl.Black)
}

/* LoadScript Phase */

type loadScript struct {
	emptyPhase

	scriptPath string
	sceneCache *scenebuilder.SceneCache

	abort bool
}

func LoadScriptPhase(sceneCache *scenebuilder.SceneCache, scriptPath string) AppPhase {
	return &loadScript{
		scriptPath: scriptPath,
		sceneCache: sceneCache,
		abort:      false,
	}
}

func (l *loadScript) Start() error {
	l.abort = false

	err := script.LoadSceneScript(l.sceneCache, l.scriptPath)
	if err != nil {
		fmt.Println(err)
		l.abort = true
	}

	return nil
}

func (l *loadScript) Think() (AppPhaseId, error) {
	if l.sceneCache.Full() || l.abort {
		return Preview, nil
	}

	return LoadScript, nil
}

func (l *loadScript) End() error {
	config := l.sceneCache.Config()
	rl.SetTargetFPS(int32(config.FrameSpeed))
	rl.SetWindowSize(int(config.ImageWidth), int(config.ImageHeight))

	return nil
}

func (l *loadScript) Draw() {
	rl.DrawText("Caching scenes...", 10, 10, 30, rl.Black)
}
