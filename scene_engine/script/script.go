package script

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
	"github.com/edison-moreland/SceneEngine/scene_engine/script/libraries"
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
	moduleMap.AddBuiltinModule("vec3", libraries.Vec3Module)
	moduleMap.AddBuiltinModule("color", libraries.ColorModule)
	moduleMap.AddBuiltinModule("shape", libraries.ShapeModule)
	moduleMap.AddBuiltinModule("material", libraries.MaterialModule)
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
						X: 4,
						Y: 0,
						Z: 0,
					},
				},
			}

			return &tengo.Map{Value: map[string]tengo.Object{
				"frame":   &tengo.Int{Value: request.frame},
				"seconds": &tengo.Float{Value: request.seconds},
				"scene_gen": &tengo.Map{Value: map[string]tengo.Object{

					"Object": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
						if len(args) != 2 {
							return nil, tengo.ErrWrongNumArguments
						}

						shape, ok := args[0].(*libraries.Shape)
						if !ok {
							return nil, tengo.ErrInvalidArgumentType{Name: "shape"}
						}

						material, ok := args[1].(*libraries.Material)
						if !ok {
							return nil, tengo.ErrInvalidArgumentType{Name: "material"}
						}

						scene.Objects = append(scene.Objects, messages.Object{
							Material: material.Material,
							Shape:    shape.Shape,
						})

						return nil, nil
					}},
					"Camera": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
						argCount := len(args)
						if argCount < 2 || argCount > 4 {
							return nil, tengo.ErrWrongNumArguments
						}

						var ok bool

						// First two are always look_from and look_at
						lookFrom, ok := args[0].(*libraries.Vec3)
						if !ok {
							return nil, tengo.ErrInvalidArgumentType{Name: "look_from"}
						}

						lookAt, ok := args[1].(*libraries.Vec3)
						if !ok {
							return nil, tengo.ErrInvalidArgumentType{Name: "look_at"}
						}

						fov := float64(90)
						aperture := 0.1
						switch argCount {
						case 4:
							aperture, ok = tengo.ToFloat64(args[3])
							if !ok {
								return nil, tengo.ErrInvalidArgumentType{Name: "aperture"}
							}

							fallthrough
						case 3:
							fov, ok = tengo.ToFloat64(args[2])
							if !ok {
								return nil, tengo.ErrInvalidArgumentType{Name: "fov"}
							}

						}

						scene.Camera = messages.Camera{
							Aperture: aperture,
							Fov:      fov,
							LookAt:   lookAt.Position(),
							LookFrom: lookFrom.Position(),
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
	config.FrameCount = 1
	config.Depth = 50
	config.Samples = 50

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
		case "frame_count":
			frameCount, ok := tengo.ToInt(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.FrameCount = uint64(frameCount)
		case "frame_speed":
			frameSpeed, ok := tengo.ToInt(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.FrameSpeed = uint64(frameSpeed)
		default:
			return config, fmt.Errorf("%w: %s", ErrUnknownConfigValue, key)
		}
	}

	config.ImageHeight = uint64(float64(config.ImageWidth) / config.AspectRatio)

	return config, nil
}

// Opaque is a Tengo object holding a value invisible to the script runtime
type Opaque struct {
	tengo.ObjectImpl
	Value any
}

func (o *Opaque) TypeName() string {
	return "Opaque"
}

func (o *Opaque) String() string {
	return fmt.Sprintf("Opaque(%T)", o.Value)
}
