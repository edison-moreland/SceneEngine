package core

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/edison-moreland/SceneEngine/src/core/messages"
	"github.com/edison-moreland/SceneEngine/submsg/runtime/go"
)

// RenderCore provides an interface to the external rendering engine
type RenderCore struct {
	client *messages.CoreClient
	info   string

	ready chan bool // notify when RenderCore is ready
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

	return &core, nil
}

// Respond to messages from the core

func (r *RenderCore) CoreReady(_ io.Reader) error {
	r.ready <- true
	return nil
}

func (r *RenderCore) CoreInfo(body io.Reader) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	var coreInfo messages.MsgCoreInfo
	err = msgpack.Unmarshal(b, &coreInfo)
	if err != nil {
		return err
	}

	r.info = coreInfo.Version
	r.ready <- true

	return nil
}

// WaitForReady will wait until the Core is ready for another command
func (r *RenderCore) WaitForReady() {
	<-r.ready
}

// Info returns the core's info string
func (r *RenderCore) Info() string {
	if r.info == "" {
		r.client.Info(0, nil)

		<-r.ready
	}
	return r.info
}
