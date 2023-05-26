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

	aspectRatio := float64(16.0 / 9.0)
	width := uint64(400)
	height := uint64(float64(width) / aspectRatio)
	renderCore.SetConfig(messages.Config{
		AspectRatio: aspectRatio,
		ImageWidth:  width,
		ImageHeight: height,
		Samples:     50,
		Depth:       50,
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

		rl.ClearBackground(rl.Blue)

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

func defaultScene() messages.Scene {
	objects := []messages.Object{
		{
			Material: messages.MaterialFrom(
				messages.Lambert{Albedo: messages.Color{
					R: 204,
					G: 204,
					B: 0,
				}},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: 0,
						Y: -100.5,
						Z: -1,
					},
					Radius: 100,
				},
			),
		},
		{
			Material: messages.MaterialFrom(
				messages.Lambert{
					Albedo: messages.Color{
						R: 178,
						G: 76,
						B: 76,
					},
				},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: 0,
						Y: 0,
						Z: -1,
					},
					Radius: 0.5,
				},
			),
		},
		{
			Material: messages.MaterialFrom(
				messages.Metal{
					Albedo: messages.Color{
						R: 204,
						G: 204,
						B: 204,
					},
					Scatter: 1.0,
				},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: -1,
						Y: 0,
						Z: -1,
					},
					Radius: 0.5,
				},
			),
		},
		{
			Material: messages.MaterialFrom(
				messages.Metal{
					Albedo: messages.Color{
						R: 204,
						G: 153,
						B: 51,
					},
					Scatter: 1.0,
				},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: 1,
						Y: 0,
						Z: -1,
					},
					Radius: 0.5,
				},
			),
		},
	}

	return messages.Scene{
		Objects: objects,
		Camera: messages.Camera{Origin: messages.Position{
			X: 0,
			Y: 0,
			Z: 0,
		}},
	}
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
