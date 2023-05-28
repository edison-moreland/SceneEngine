package script

import (
	"errors"
	"fmt"
	"os"

	"github.com/d5/tengo/v2"

	"github.com/edison-moreland/SceneEngine/src/core/messages"
)

type SceneScript struct {
	compiled *tengo.Compiled
}

func LoadSceneScript(path string) (*SceneScript, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	script := tengo.NewScript(source)

	compiled, err := script.Run()
	if err != nil {
		return nil, err
	}

	return &SceneScript{
		compiled: compiled,
	}, nil
}

var (
	ErrConfigIsUndefined  = errors.New("config is undefined")
	ErrConfigIsNotMap     = errors.New("config is not a map")
	ErrUnknownConfigValue = errors.New("unknown config value")
)

func (s *SceneScript) Config() (messages.Config, error) {
	var config messages.Config

	configVar := s.compiled.Get("config")
	if configVar.IsUndefined() {
		return config, ErrConfigIsUndefined
	}

	configMap := configVar.Map()
	if configMap == nil {
		return config, ErrConfigIsNotMap
	}

	for key, val := range configMap {
		switch key {
		case "depth":
			config.Depth = uint64(val.(int64))
		case "samples":
			config.Samples = uint64(val.(int64))
		case "image_width":
			config.ImageWidth = uint64(val.(int64))
		case "aspect_ratio":
			config.AspectRatio = val.(float64)
		default:
			return config, fmt.Errorf("%w: %s", ErrUnknownConfigValue, key)
		}
	}

	config.ImageHeight = uint64(float64(config.ImageWidth) / config.AspectRatio)

	return config, nil
}
