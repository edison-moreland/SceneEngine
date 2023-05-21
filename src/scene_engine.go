package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/edison-moreland/SceneEngine/src/core"
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
	fmt.Println("Starting Scene Engine")
	coreCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	renderCore, err := core.Start(coreCtx, corePath)
	if err != nil {
		fmt.Println("Could not start core!")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	renderCore.WaitForReady()
	fmt.Printf("Using core: %s \n", renderCore.Info())

	fmt.Println("Starting render.")
	for pixel := range renderCore.StartRender() {
		fmt.Printf("Received a pixel! %#v \n", pixel)
	}

	fmt.Println("Done. goodbye!")
}
