package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"go.uber.org/zap"

	"github.com/edison-moreland/SceneEngine/scene_engine/core"
	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
	"github.com/edison-moreland/SceneEngine/scene_engine/script"
)

type EnginePhase int

const (
	Idle EnginePhase = iota
	Rendering
)

var (
	corePath   string
	scriptPath string
)

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
	config, requestScene, err := script.LoadSceneScript(engineCtx, scriptPath)
	if err != nil {
		logger.Fatal("Could not load script!", zap.Error(err))
	}

	exportDir, err := prepareFS(scriptPath)
	if err != nil {
		logger.Fatal("Could not prepare export dir!", zap.Error(err))
	}

	logger.Info("Starting render core")
	renderCore, err := core.Start(engineCtx, corePath)
	if err != nil {
		logger.Fatal("Could not start core!", zap.Error(err))
	}
	renderCore.WaitForReady()
	renderCore.SetConfig(config)
	renderCore.WaitForReady()
	logger.Info("Render core ready!", zap.String("version", renderCore.Info()))

	rl.InitWindow(int32(config.ImageWidth), int32(config.ImageHeight), "SceneEngine")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	currentFrame := uint64(1)

	startRender := func(frame uint64) <-chan []messages.Pixel {
		logger.Info("Rendering frame", zap.Uint64("frame", currentFrame))

		return renderCore.StartRender(requestScene(uint32(frame), float64(frame)*(1.0/float64(config.FrameSpeed))))
	}

	target := newRenderTarget(
		config.ImageWidth,
		config.ImageHeight,
		startRender(currentFrame))

	phase := Rendering

drawLoop:
	for !rl.WindowShouldClose() {
		switch phase {
		case Rendering:
			target.RenderBuffer()
			if target.done {
				phase = Idle
			}

		case Idle:
			if target != nil {
				target.Export(exportDir, currentFrame)

				currentFrame += 1
				if currentFrame > config.FrameCount {
					break drawLoop
				}

				phase = Rendering
				renderCore.WaitForReady() // This should return immediately
				target.Reset(startRender(currentFrame))
			}
		}

		rl.BeginDrawing()

		switch phase {
		case Rendering:
			rl.DrawTextureV(target.Texture, rl.Vector2Zero(), rl.White)
		}

		rl.EndDrawing()
	}
	logger.Info("Done!")

	if config.FrameSpeed > 0 {
		logger.Info("Exporting video")
		cmd, err := ffmpegEncodeVideo(engineCtx, config, exportDir)
		if err != nil {
			log.Fatal("Could not start ffmpeg to export video", zap.Error(err))
		}
		err = cmd.Wait()
		if err != nil {
			log.Fatal("Err exporting video", zap.Error(err))
		}
	}
}

type renderTarget struct {
	sync.Mutex
	rl.RenderTexture2D
	buffer []messages.Pixel
	done   bool
}

func newRenderTarget(width uint64, height uint64, pixels <-chan []messages.Pixel) *renderTarget {
	var r renderTarget
	r.buffer = make([]messages.Pixel, 0, 255) // TODO: How big should the buffer be to start?
	r.RenderTexture2D = rl.LoadRenderTexture(int32(width), int32(height))
	r.done = true

	r.Reset(pixels)
	//go func() {
	//	for batch := range pixels {
	//		r.Lock()
	//		r.buffer = append(r.buffer, batch...)
	//		r.Unlock()
	//	}
	//	r.done = true
	//}()

	return &r
}

// RenderBuffer is called in the main thread to render buffered pixels to the texture
func (r *renderTarget) RenderBuffer() {
	r.Lock()
	if len(r.buffer) == 0 {
		r.Unlock()
		return
	}
	defer r.Unlock()

	rl.BeginTextureMode(r.RenderTexture2D)
	defer rl.EndTextureMode()

	for _, pixel := range r.buffer {
		rl.DrawPixel(
			int32(pixel.X),
			int32(pixel.Y),
			rl.NewColor(
				pixel.Color.R,
				pixel.Color.G,
				pixel.Color.B,
				255,
			),
		)
	}

	r.buffer = r.buffer[:0]
}

func (r *renderTarget) Export(dir string, frame uint64) {
	exportName := path.Join(dir, fmt.Sprintf("frame-%d.png", frame))

	r.Lock()
	image := rl.LoadImageFromTexture(r.Texture)
	rl.ExportImage(*image, exportName)
	r.Unlock()
}

func (r *renderTarget) Reset(pixels <-chan []messages.Pixel) {
	r.Lock()

	if !r.done {
		panic("can't reset if not done")
	}
	r.done = false

	r.Unlock()

	go func() {
		for batch := range pixels {
			r.Lock()
			r.buffer = append(r.buffer, batch...)
			r.Unlock()
		}
		r.done = true
	}()
}

func (r *renderTarget) Close() {
	r.Lock()
	rl.UnloadRenderTexture(r.RenderTexture2D)
	r.Unlock()
}

// prepareFS will make sure an empty folder exists for output, removing any old files
func prepareFS(scriptPath string) (string, error) {
	baseDir, scritName := path.Split(scriptPath)
	name := strings.SplitN(scritName, ".", 2)[0]

	exportDir := path.Join(baseDir, name)

	if err := os.RemoveAll(exportDir); err != nil {
		return "", err
	}

	return exportDir, os.Mkdir(exportDir, os.FileMode(0777))
}
