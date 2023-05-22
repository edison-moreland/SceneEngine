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
)

var corePath string

func init() {
	flag.StringVar(&corePath, "core", "", "Path to rendercore")
	flag.Parse()

	if corePath == "" {
		fmt.Println("parameter `core` is required.")
		os.Exit(1)
	}
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("Welcome to SceneEngine!")

	logger.Info("Starting render core")
	coreCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	renderCore, err := core.Start(coreCtx, corePath)
	if err != nil {
		logger.Fatal("Could not start core!", zap.Error(err))
	}

	renderCore.WaitForReady()
	logger.Info("Render core ready!", zap.String("version", renderCore.Info()))

	width := uint64(800)
	height := uint64(450)
	renderCore.SetConfig(messages.Config{
		ImageWidth:  width,
		ImageHeight: height,
	})
	renderCore.WaitForReady()
	logger.Info("Set config")

	rl.InitWindow(int32(width), int32(height), "SceneEngine")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	logger.Info("Starting render")
	target := NewRenderTarget(width, height, renderCore.StartRender())

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		for x := uint64(0); x < width; x++ {
			for y := uint64(0); y < height; y++ {
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
			r.image[(p.Y*height)+p.X] = rl.ColorFromNormalized(rl.Vector4{
				X: float32(p.Color.R),
				Y: float32(p.Color.G),
				Z: float32(p.Color.B),
				W: 1,
			})
		}
	}()

	return &r
}

func (r *renderTarget) Pixel(x, y uint64) (rl.Vector2, rl.Color) {
	color := r.image[(y*r.height)+x]

	return rl.Vector2{
		X: float32(x),
		Y: float32(y),
	}, color
}
