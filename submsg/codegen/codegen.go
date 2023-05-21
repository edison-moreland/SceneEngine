package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	descriptorPath string
	goOutputFile   string
	goPackage      string
	ponyOutputFile string
)

func init() {
	flag.StringVar(&descriptorPath, "descriptor", "", "Path to submsg.yaml")
	flag.StringVar(&goOutputFile, "go-output", "", "File path to output generated go code")
	flag.StringVar(&goPackage, "go-package", "", "Package name to use when generating go code")
	flag.StringVar(&ponyOutputFile, "pony-output", "", "File path to output generated pony code")
	flag.Parse()

	if descriptorPath == "" {
		fmt.Println("descriptor required!")
		os.Exit(1)
	}
	if goOutputFile == "" {
		fmt.Println("go-output required!")
		os.Exit(1)
	}
	if goPackage == "" {
		fmt.Println("go-package required!")
		os.Exit(1)
	}
	if ponyOutputFile == "" {
		fmt.Println("pony-output required!")
		os.Exit(1)
	}
}

type codegen interface {
	MsgIds(string, []MsgDesc)
	Server(string, []MsgDesc)
	Client(string, []MsgDesc)
	Done() error
}

func main() {
	desc, err := loadDescriptor(descriptorPath)
	if err != nil {
		panic(err)
	}

	err = generate(desc, newGoGen(), newPonyGen())
	if err != nil {
		panic(err)
	}
}

func generate(desc SubMsgDesc, generators ...codegen) error {
	for _, g := range generators {
		generatePackage(desc.Messages, g)

		if err := g.Done(); err != nil {
			return err
		}
	}
	return nil
}

func generatePackage(desc MessagesDesc, gen codegen) {
	gen.MsgIds("engine", desc.Engine)
	gen.Server("engine", desc.Engine)
	gen.Client("engine", desc.Engine)
	gen.MsgIds("core", desc.Core)
	gen.Server("core", desc.Core)
	gen.Client("core", desc.Core)
}
