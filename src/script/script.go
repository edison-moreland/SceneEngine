package script

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"

	"github.com/edison-moreland/SceneEngine/src/core/messages"
)

var (
	ErrConfigIsNotMap           = errors.New("config is not a map")
	ErrUnknownConfigValue       = errors.New("unknown config value")
	ErrConfigValueIncorrectType = errors.New("config value has the wrong type")
)

//go:embed runtime.tengo
var runtimeSource []byte

type GenerateScene func(frame uint32, seconds float64) messages.Scene

type sceneRequest struct {
	frame    int64
	seconds  float64
	response chan<- messages.Scene
}

func LoadSceneScript(ctx context.Context, sceneScript string) (messages.Config, GenerateScene, error) {
	request := make(chan sceneRequest)

	config, err := startScript(ctx, sceneScript, request)
	if err != nil {
		close(request)
		return config, nil, err
	}

	return config, func(frame uint32, seconds float64) messages.Scene {
		responseChan := make(chan messages.Scene)
		defer close(responseChan)

		request <- sceneRequest{
			frame:    int64(frame),
			seconds:  seconds,
			response: responseChan,
		}

		return <-responseChan
	}, nil
}

func startScript(ctx context.Context, sceneScript string, requests chan sceneRequest) (messages.Config, error) {
	var config messages.Config

	source, err := os.ReadFile(sceneScript)
	if err != nil {
		return config, err
	}

	configReturn := make(chan messages.Config)
	defer close(configReturn)

	moduleMap := tengo.NewModuleMap()
	moduleMap.AddSourceModule("userscript", source)
	moduleMap.AddBuiltinModule("fmt", stdlib.BuiltinModules["fmt"])
	moduleMap.AddBuiltinModule("math", stdlib.BuiltinModules["math"])
	moduleMap.AddBuiltinModule("runtime", map[string]tengo.Object{
		"config": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			// Runtime is giving us the config defined by the userscript
			if configReturn == nil {
				panic("config can only be called once")
			}

			if len(args) != 1 {
				panic("config expects more than one argument")
			}

			configMap, ok := args[0].(*tengo.Map)
			if !ok {
				return nil, ErrConfigIsNotMap
			}

			userConfig, err := getConfig(configMap)
			if err != nil {
				return nil, err
			}

			configReturn <- userConfig

			return nil, nil
		}},
		"next": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			// Runtime is asking for the next request
			// When a request is ready, it returns an object with:
			//   - A callback to call when request is done
			//   - An object for creating scene objects
			//   - The two scene arguments (frame, seconds)
			request := <-requests

			// This scene is populated by the scene_gen object
			scene := messages.Scene{
				Camera: messages.Camera{
					Aperture: 0.1,
					Fov:      90,
					LookAt: messages.Position{
						X: 0,
						Y: 0,
						Z: 0,
					},
					LookFrom: messages.Position{
						X: 5,
						Y: 5,
						Z: 5,
					},
				},
			}

			return &tengo.Map{Value: map[string]tengo.Object{
				"frame":   &tengo.Int{Value: request.frame},
				"seconds": &tengo.Float{Value: request.seconds},
				"scene_gen": &tengo.Map{Value: map[string]tengo.Object{
					"Spheres": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
						if len(args) != 1 {
							return nil, fmt.Errorf("expected 1 argument")
						}

						sphereCount, ok := tengo.ToInt(args[0])
						if !ok {
							return nil, fmt.Errorf("expected first argument to be int")
						}

						for i := 0; i < sphereCount; i++ {
							scene.Objects = append(scene.Objects, messages.Object{
								Material: messages.MaterialFrom(messages.Lambert{Albedo: messages.Color{
									B: 150,
									G: 150,
									R: 150,
								}}),
								Shape: messages.ShapeFrom(messages.Sphere{
									Origin: messages.Position{
										X: float64(i),
										Y: 0,
										Z: 0,
									},
									Radius: 1,
								}),
							})
						}

						return nil, nil
					}},
				}},
				"done": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
					request.response <- scene

					return nil, nil
				}},
			}}, nil

		}},
	})

	runtime := tengo.NewScript(runtimeSource)
	runtime.SetImports(moduleMap)

	go func() {
		defer close(requests)

		_, err := runtime.RunContext(ctx)
		if err != nil {
			panic(err)
		}
	}()

	return <-configReturn, nil
}

func getConfig(o *tengo.Map) (messages.Config, error) {
	var config messages.Config

	for key, val := range o.Value {
		switch key {
		case "depth":
			depth, ok := tengo.ToInt(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.Depth = uint64(depth)
		case "samples":
			samples, ok := tengo.ToInt(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.Samples = uint64(samples)
		case "image_width":
			imageWidth, ok := tengo.ToInt(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.ImageWidth = uint64(imageWidth)
		case "aspect_ratio":
			aspectRatio, ok := tengo.ToFloat64(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.AspectRatio = aspectRatio
		default:
			return config, fmt.Errorf("%w: %s", ErrUnknownConfigValue, key)
		}
	}

	config.ImageHeight = uint64(float64(config.ImageWidth) / config.AspectRatio)

	return config, nil
}
