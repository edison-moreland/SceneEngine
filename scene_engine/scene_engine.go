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

	logger.Info("Loading scene script")
	config, scene, err := script.LoadSceneScript(engineCtx, scriptPath)
	if err != nil {
		logger.Fatal("Could not load script!", zap.Error(err))
	}

	logger.Info("Starting render core")
	renderCore, err := core.Start(engineCtx, corePath)
	if err != nil {
		logger.Fatal("Could not start core!", zap.Error(err))
	}
	logger.Info("Render core ready!", zap.String("version", renderCore.Info()))
	renderCore.SetConfig(config)

	rl.InitWindow(int32(config.ImageWidth), int32(config.ImageHeight), "SceneEngine")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	// TODO:
	// - LoadScript <- starting phase
	//   - populates scene cache
	// 	 - fancy loading bar
	// - EncodeVideo
	//   - another fancy loading bar
	//   - fix `<defunct>` pids left behind by encode
	err = RunPhases(logger, Idle, map[AppPhaseId]AppPhase{
		Idle:   IdlePhase(),
		Render: RenderPhase(renderCore, scene, config, exportDir(scriptPath)),
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
	Render
)

type AppPhase interface {
	Think() (next AppPhaseId, err error)
	Draw()
	Complete() error
}

func RunPhases(logger *zap.Logger, startPhase AppPhaseId, phases map[AppPhaseId]AppPhase) error {
	current := startPhase

	defer func() {
		for _, phase := range phases {
			if err := phase.Complete(); err != nil {
				logger.Error("Err completing phase", zap.Error(err))
			}
		}
	}()

	for !rl.WindowShouldClose() {
		next, err := phases[current].Think()
		if err != nil {
			return err
		}

		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)
		phases[current].Draw()
		rl.EndDrawing()

		current = next
	}

	return nil
}

/* Idle Phase */

type idle struct {
	renderButton Button
}

func IdlePhase() AppPhase {
	return &idle{
		renderButton: NewButton("Render", rl.Gray, 5, 5),
	}
}

func (i *idle) Think() (AppPhaseId, error) {
	if i.renderButton.Down() {
		return Render, nil
	}

	return Idle, nil
}

func (i *idle) Draw() {
	i.renderButton.Draw()
}

func (i *idle) Complete() error {
	return nil
}

/* Render phase */

type render struct {
	core       *core.RenderCore
	target     *renderTarget
	scene      script.GenerateScene
	config     messages.Config
	exportPath string

	exportCmd *exec.Cmd

	renderActive bool

	currentFrame   uint64
	lastFrameStart time.Time
	frameElapsed   rollingAverage
}

func RenderPhase(core *core.RenderCore, scene script.GenerateScene, config messages.Config, exportDir string) AppPhase {
	r := render{
		core:       core,
		scene:      scene,
		config:     config,
		exportPath: exportDir,

		renderActive: false,

		currentFrame: 0,
		frameElapsed: rollingAverage{sampleSize: 10},
	}

	r.target = newRenderTarget(config.ImageWidth, config.ImageHeight)

	return &r
}

func (r *render) prepareExportDir() error {
	if err := os.RemoveAll(r.exportPath); err != nil {
		return err
	}

	return os.Mkdir(r.exportPath, os.FileMode(0777))
}

func (r *render) startFrame() {
	scene := r.scene(
		uint32(r.currentFrame),
		float64(r.currentFrame)*(1.0/float64(r.config.FrameSpeed)),
	)

	r.lastFrameStart = time.Now()
	r.core.WaitForReady() // This should return immediately
	r.target.Reset(r.core.StartRender(scene))
}

func (r *render) startExport() {
	if r.exportCmd != nil {
		if r.exportCmd.ProcessState != nil {
			// TODO: Double check that this check actually works
			panic("refusing to start exporting when already exporting")
		}
	}

	inputImages := path.Join(r.exportPath, "frame-%d.png")
	outputVideo := path.Clean(path.Join(r.exportPath, "..", path.Base(r.exportPath)+".mp4"))

	r.exportCmd = exec.Command("ffmpeg",
		"-y",
		"-framerate", strconv.Itoa(int(r.config.FrameSpeed)),
		"-i", inputImages,
		"-pix_fmt", "yuv420p",
		"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2",
		"-vcodec", "h264",
		"-acodec", "aac",
		outputVideo,
	)
	r.exportCmd.Stderr = os.Stderr
	r.exportCmd.Stdout = os.Stdout

	if err := r.exportCmd.Start(); err != nil {
		panic(err)
	}
}

func (r *render) Think() (AppPhaseId, error) {
	if !r.renderActive {
		if r.currentFrame == 0 {
			if err := r.prepareExportDir(); err != nil {
				return Idle, err
			}
		}

		r.currentFrame += 1
		if r.currentFrame > r.config.FrameCount {
			r.startExport()
			r.currentFrame = 0
			r.frameElapsed.Reset()
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

func (r *render) Draw() {
	rl.DrawTextureV(r.target.Texture, rl.Vector2Zero(), rl.White)

	DrawInfo(0, "frame", fmt.Sprintf("%d/%d", r.currentFrame, r.config.FrameCount))
	if r.frameElapsed.HasSamples() {
		averageFrameTime := r.frameElapsed.Average()
		DrawInfo(1, "avg_frame_time", averageFrameTime.Round(time.Millisecond).String())
		DrawInfo(2, "done_in", (time.Duration(r.config.FrameCount-r.currentFrame) * averageFrameTime).Round(time.Millisecond).String())
	}
}

func (r *render) Complete() error {
	if r.exportCmd != nil {
		return r.exportCmd.Wait()
	}

	return nil
}
