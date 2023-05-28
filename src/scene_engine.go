package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
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

	aspectRatio := float64(3.0 / 2.0)
	width := uint64(1200)
	height := uint64(float64(width) / aspectRatio)
	renderCore.SetConfig(messages.Config{
		AspectRatio: aspectRatio,
		ImageWidth:  width,
		ImageHeight: height,
		Samples:     500,
		Depth:       50,
	})
	renderCore.WaitForReady()
	logger.Info("Set config")

	rl.InitWindow(int32(width), int32(height), "SceneEngine")

	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	logger.Info("Starting render")
	target := NewRenderTarget(width, height, renderCore.StartRender(defaultScene()))

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
					R: 127,
					G: 127,
					B: 127,
				}},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: 0,
						Y: -1000,
						Z: 0,
					},
					Radius: 1000,
				},
			),
		},
		{
			Material: messages.MaterialFrom(
				messages.Dielectric{
					IndexOfRefraction: 1.5,
				},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: 0,
						Y: 1,
						Z: 0,
					},
					Radius: 1,
				},
			),
		},
		{
			Material: messages.MaterialFrom(
				messages.Lambert{
					Albedo: messages.Color{
						R: 102,
						G: 51,
						B: 25,
					},
				},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: -4,
						Y: 1,
						Z: 0,
					},
					Radius: 1,
				},
			),
		},
		{
			Material: messages.MaterialFrom(
				messages.Metal{
					Albedo: messages.Color{
						R: 178,
						G: 153,
						B: 127,
					},
				},
			),
			Shape: messages.ShapeFrom(
				messages.Sphere{
					Origin: messages.Position{
						X: 4,
						Y: 1,
						Z: 0,
					},
					Radius: 1,
				},
			),
		},
	}

	for a := -11; a < 11; a += 10 {
		for b := -11; b < 11; b += 10 {
			center := rl.NewVector3(
				float32(a)+(0.9*rand.Float32()),
				0.2,
				float32(b)+(0.9*rand.Float32()),
			)

			if rl.Vector3Length(rl.Vector3Subtract(center, rl.NewVector3(4, 0.2, 0))) > 0.9 {
				materialChoice := rand.Float32()
				var material messages.Material
				if materialChoice < 0.8 {
					material = messages.MaterialFrom(messages.Lambert{
						Albedo: messages.Color{
							R: uint8(rand.Int()),
							G: uint8(rand.Int()),
							B: uint8(rand.Int()),
						},
					})
				} else if materialChoice < 0.95 {
					material = messages.MaterialFrom(messages.Metal{
						Albedo: messages.Color{
							R: uint8(rand.Intn(125) + 125),
							G: uint8(rand.Intn(125) + 125),
							B: uint8(rand.Intn(125) + 125),
						},
						Scatter: rand.Float64() / 2,
					})
				} else {
					material = messages.MaterialFrom(messages.Dielectric{
						IndexOfRefraction: 1.5,
					})
				}

				objects = append(objects, messages.Object{
					Material: material,
					Shape: messages.ShapeFrom(messages.Sphere{
						Origin: messages.Position{
							X: float64(center.X),
							Y: float64(center.Y),
							Z: float64(center.Z),
						},
						Radius: 0.2,
					}),
				})
			}
		}
	}

	return messages.Scene{
		Objects: objects,
		Camera: messages.Camera{
			LookFrom: messages.Position{
				X: 13,
				Y: 2,
				Z: 3,
			},
			LookAt: messages.Position{
				X: 0,
				Y: 0,
				Z: 0,
			},
			Fov:      20,
			Aperture: 0.1,
		},
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
