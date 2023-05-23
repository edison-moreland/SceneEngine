package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func loadDescriptor(path string) (SubMsgDesc, error) {
	var desc SubMsgDesc

	content, err := os.ReadFile(path)
	if err != nil {
		return desc, err
	}

	err = yaml.Unmarshal(content, &desc)
	if err != nil {
		return desc, err
	}

	return desc, nil
}

// SubMsgDesc - Defines the structure of submsg.yaml
type SubMsgDesc struct {
	Types    []TypeDesc
	Messages map[string][]MsgDesc
}

type MsgDesc struct {
	Name        string
	Description string
}

type TypeDesc struct {
	Name   string
	Type   string
	Fields map[string]string
}
