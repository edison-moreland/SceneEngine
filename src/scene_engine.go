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
	target := NewRenderTarget(
		config.ImageWidth,
		config.ImageHeight,
		renderCore.StartRender(scene),
	)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.Blue)

		for x := uint64(0); x < config.ImageWidth; x++ {
			for y := uint64(0); y < config.ImageHeight; y++ {
				rl.DrawPixelV(target.Pixel(x, y))
			}
		}

		rl.EndDrawing()
	}
}

type renderTarget struct {
	sync.RWMutex
	image  []rl.Color
	height uint64
}

func NewRenderTarget(width uint64, height uint64, pixels <-chan messages.Pixel) *renderTarget {
	var r renderTarget
	r.image = make([]rl.Color, width*height)
	r.height = height

	go func() {
		for p := range pixels {
			// TODO: Make sure pixels are laid out in the same order that they're read back
			// TODO: Adjust submsg types to use float32

			r.image[(p.X*height)+p.Y] = rl.NewColor(
				p.Color.R,
				p.Color.G,
				p.Color.B,
				255,
			)
		}
	}()

	return &r
}

func (r *renderTarget) Pixel(x, y uint64) (rl.Vector2, rl.Color) {
	color := r.image[(x*r.height)+y]

	return rl.Vector2{
		X: float32(x),
		Y: float32(y),
	}, color
}
