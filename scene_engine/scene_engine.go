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

	var sceneCache script.SceneCache
	err = RunPhases(logger, LoadScript, map[AppPhaseId]AppPhase{
		Idle:       IdlePhase(scriptFileWatcher),
		LoadScript: LoadScriptPhase(&sceneCache, scriptPath),
		Render:     RenderPhase(renderCore, &sceneCache, exportDir(scriptPath)),
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
	Idle AppPhaseId = iota
	LoadScript
	Render
	Encode
)

type AppPhase interface {
	Think() (next AppPhaseId, err error)
	Draw()

	Start() error // Called when transitioning to this phase
	End() error   // Called when transitioning to a different phase

	Shutdown() error // Ran when app is shutting down
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

/* Idle Phase */

type idle struct {
	emptyPhase
	script       *FileWatcher
	renderButton Button
}

func IdlePhase(script *FileWatcher) AppPhase {
	return &idle{
		script:       script,
		renderButton: NewButton("Render", rl.Gray, 5, 5),
	}
}

func (i *idle) Think() (AppPhaseId, error) {
	if i.renderButton.Down() {
		return Render, nil
	}

	if i.script.HasChanged() {
		return LoadScript, nil
	}

	return Idle, nil
}

func (i *idle) End() error {
	i.script.ClearChange()

	return nil
}

func (i *idle) Draw() {
	i.renderButton.Draw()
}

/* Render Phase */

type render struct {
	emptyPhase

	core       *core.RenderCore
	target     *renderTarget
	sceneCache *script.SceneCache
	exportPath string

	renderActive bool

	currentFrame   uint64
	lastFrameStart time.Time
	frameElapsed   rollingAverage
}

func RenderPhase(core *core.RenderCore, sceneCache *script.SceneCache, exportDir string) AppPhase {
	r := render{
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
	rl.SetWindowSize(int(config.ImageWidth), int(config.ImageHeight))
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
			return Idle, nil
		}

		r.startFrame()
		r.renderActive = true
	}

	r.target.RenderBufferToTexture()
	if r.target.done {
		r.frameElapsed.Sample(time.Since(r.lastFrameStart))
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
	DrawInfo(0, "frame", fmt.Sprintf("%d/%d", r.currentFrame, config.FrameCount))
	if r.frameElapsed.HasSamples() {
		averageFrameTime := r.frameElapsed.Average()
		DrawInfo(1, "avg_frame_time", averageFrameTime.Round(time.Millisecond).String())
		DrawInfo(2, "done_in", (time.Duration(config.FrameCount-r.currentFrame) * averageFrameTime).Round(time.Millisecond).String())
	}
}

func (r *render) Shutdown() error {
	return nil
}

/* Encode Phase */

type encode struct {
	emptyPhase

	sceneCache *script.SceneCache

	encodeCmd    *exec.Cmd
	encodeActive bool
	exportDir    string
}

func EncodePhase(sceneCache *script.SceneCache, exportDir string) AppPhase {
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

	return Idle, nil
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
	sceneCache *script.SceneCache

	abort bool
}

func LoadScriptPhase(sceneCache *script.SceneCache, scriptPath string) AppPhase {
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
		return Idle, nil
	}

	return LoadScript, nil
}

func (l *loadScript) Draw() {
	rl.DrawText("Caching scenes...", 10, 10, 30, rl.Black)
}
