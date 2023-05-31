package main

import (
	"fmt"
	"path"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

type renderTarget struct {
	sync.Mutex
	rl.RenderTexture2D
	buffer []messages.Pixel
	done   bool
}

func newRenderTarget(width uint64, height uint64) *renderTarget {
	var r renderTarget
	r.buffer = make([]messages.Pixel, 0, 255) // TODO: How big should the buffer be to start?
	r.RenderTexture2D = rl.LoadRenderTexture(int32(width), int32(height))
	r.done = true

	return &r
}

// RenderBufferToTexture is called in the main thread to render buffered pixels to the texture
func (r *renderTarget) RenderBufferToTexture() {
	r.Lock()
	defer r.Unlock()

	if len(r.buffer) == 0 {
		return
	}

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

	rl.BeginTextureMode(r.RenderTexture2D)
	rl.ClearBackground(backgroundColor)
	rl.EndTextureMode()

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
