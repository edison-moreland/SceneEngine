package core

import (
	"context"
	"fmt"
	"os"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
	"github.com/edison-moreland/SceneEngine/submsg/runtime/go"
)

// RenderCore provides an interface to the external rendering engine
type RenderCore struct {
	client *messages.CoreClient
	info   string

	ready chan bool // notify when RenderCore is ready

	pixelsOut chan []messages.Pixel
}

func Start(ctx context.Context, path string) (*RenderCore, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("%w: Failed to stat core at %s", err, path)
	}

	core := RenderCore{
		ready: make(chan bool),
	}

	sendMsg, err := submsg.Start(ctx, path, messages.EngineRouter(&core))
	if err != nil {
		return nil, err
	}
	core.client = messages.NewCoreClient(sendMsg)
	core.WaitForReady()

	return &core, nil
}

// Respond to messages from the core

func (r *RenderCore) Test(body []byte) error {
	var material messages.Material
	err := msgpack.Unmarshal(body, &material)
	if err != nil {
		return err
	}

	fmt.Printf("%#v \n", material)

	return nil
}

func (r *RenderCore) CoreReady(_ []byte) error {
	// If this channel isn't nil, that means a render is in progress
	// The core will send the ready message when a render concludes
	// To notify the downstream that the render is complete, we need to close the channel
	// We should introduce explicit "phases", eg: "starting", "ready", "rendering"
	if r.pixelsOut != nil {
		close(r.pixelsOut)
		r.pixelsOut = nil
	}

	r.ready <- true
	return nil
}

func (r *RenderCore) CoreInfo(body []byte) error {
	var coreInfo messages.MsgCoreInfo
	err := msgpack.Unmarshal(body, &coreInfo)
	if err != nil {
		return err
	}

	r.info = coreInfo.Version
	r.ready <- true

	return nil
}

func (r *RenderCore) PixelBatch(body []byte) error {
	var pixels []messages.Pixel
	err := msgpack.Unmarshal(body, &pixels)
	if err != nil {
		return err
	}

	r.pixelsOut <- pixels

	return nil
}

// The actual public methods are below \/
// TODO: add another layer of separation here so the methods above aren't public

// WaitForReady will wait until the Core is ready for another command
func (r *RenderCore) WaitForReady() {
	<-r.ready
}

// Info returns the core's info string
func (r *RenderCore) Info() string {
	if r.info == "" {
		r.client.Info(nil)

		<-r.ready
	}
	return r.info
}

func (r *RenderCore) SetConfig(c messages.Config) {
	body, err := msgpack.Marshal(&c)
	if err != nil {
		panic(err)
	}

	r.client.Config(body)
}

// StartRender will start rendering the next frame
func (r *RenderCore) StartRender(scene messages.Scene) <-chan []messages.Pixel {
	if r.pixelsOut != nil {
		panic("Can't start rendering while already rendering")
	}
	r.pixelsOut = make(chan []messages.Pixel)

	body, err := msgpack.Marshal(&scene)
	if err != nil {
		panic(err)
	}

	r.client.RenderFrame(body)

	return r.pixelsOut
}
