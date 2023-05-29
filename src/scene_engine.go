package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
	"go.uber.org/zap"

	"github.com/edison-moreland/SceneEngine/src/core"
	"github.com/edison-moreland/SceneEngine/src/core/messages"
	"github.com/edison-moreland/SceneEngine/src/script"
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

	logger.Info("Starting render core")
	renderCore, err := core.Start(engineCtx, corePath)
	if err != nil {
		logger.Fatal("Could not start core!", zap.Error(err))
	}
	renderCore.WaitForReady()
	logger.Info("Render core ready!", zap.String("version", renderCore.Info()))

	renderCore.SetConfig(config)
	renderCore.WaitForReady()
	logger.Info("Set config")

	logger.Info("Calculating scene")
	scene := requestScene(1, 0)

	rl.InitWindow(int32(config.ImageWidth), int32(config.ImageHeight), "SceneEngine")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	logger.Info("Starting render")
	target := newRenderTarget(
		config.ImageWidth,
		config.ImageHeight,
		renderCore.StartRender(scene),
	)

	for !rl.WindowShouldClose() {
		target.RenderBuffer()

		rl.BeginDrawing()
		rl.DrawTexture(target.Texture, 0, 0, rl.White)
		rl.EndDrawing()
	}
}

type renderTarget struct {
	sync.Mutex
	rl.RenderTexture2D
	buffer []messages.Pixel
}

func newRenderTarget(width uint64, height uint64, pixels <-chan []messages.Pixel) *renderTarget {
	var r renderTarget
	r.buffer = make([]messages.Pixel, 0, 255) // TODO: How big should the buffer be to start?
	r.RenderTexture2D = rl.LoadRenderTexture(int32(width), int32(height))

	go func() {
		for batch := range pixels {
			r.Lock()
			r.buffer = append(r.buffer, batch...)
			r.Unlock()
		}
	}()

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

func (r *renderTarget) Close() {
	rl.UnloadRenderTexture(r.RenderTexture2D)
}
