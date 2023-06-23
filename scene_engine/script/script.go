package script

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/stdlib"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
	"github.com/edison-moreland/SceneEngine/scene_engine/script/builtins"
)

var (
	ErrUnknownConfigValue       = errors.New("unknown config value")
	ErrConfigValueIncorrectType = errors.New("config value has the wrong type")
)

//go:embed runtime.tengo
var runtimeSource []byte

func emptyScene() messages.Scene {
	return messages.Scene{
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
}

func LoadSceneScript(sceneCache *SceneCache, sceneScript string) error {
	globals := make([]tengo.Object, tengo.GlobalsSize)
	symbolTable := tengo.NewSymbolTable()

	currentScene := emptyScene()

	// Globals are only available to the runtime
	rtBegin := symbolTable.Define("rt_begin")
	globals[rtBegin.Index] = &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		fmt.Println("rt_begin")
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}

		configMap, ok := args[0].(*tengo.Map)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "config"}
		}

		userConfig, err := getConfig(configMap)
		if err != nil {
			return nil, err
		}

		sceneCache.Reset(userConfig)

		return &tengo.Map{Value: map[string]tengo.Object{
			"count":   &tengo.Int{Value: int64(userConfig.FrameCount)},
			"seconds": &tengo.Float{Value: 1 / float64(userConfig.FrameSpeed)},
		}}, nil
	}}

	rtCommitScene := symbolTable.Define("rt_commit_scene")
	globals[rtCommitScene.Index] = &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		fmt.Println("rt_commit_scene")
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}

		frame, ok := tengo.ToInt(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "frame"}
		}

		sceneCache.CacheScene(uint64(frame), currentScene)
		currentScene = emptyScene()

		return nil, nil
	}}

	// Builtins are available everywhere
	customBuiltins := []*tengo.BuiltinFunction{
		{
			Name: "scene",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				fmt.Println("scene")
				if len(args) != 2 {
					return nil, tengo.ErrWrongNumArguments
				}

				return &tengo.Map{Value: map[string]tengo.Object{
					"config": args[0],
					"scene":  args[1],
				}}, nil
			},
		},
		{
			Name: "object",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				fmt.Println("object")
				if len(args) != 2 {
					return nil, tengo.ErrWrongNumArguments
				}

				shape, ok := args[0].(*builtins.Shape)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{Name: "shape"}
				}

				material, ok := args[1].(*builtins.Material)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{Name: "material"}
				}

				currentScene.Objects = append(currentScene.Objects, messages.Object{
					Material: material.Material,
					Shape:    shape.Shape,
				})

				return nil, nil
			},
		},
		{
			Name: "camera",
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				fmt.Println("camera")
				argCount := len(args)
				if argCount < 2 || argCount > 4 {
					return nil, tengo.ErrWrongNumArguments
				}

				var ok bool

				// First two are always look_from and look_at
				lookFrom, ok := args[0].(*builtins.Vec3)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{Name: "look_from"}
				}

				lookAt, ok := args[1].(*builtins.Vec3)
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

				currentScene.Camera = messages.Camera{
					Aperture: aperture,
					Fov:      fov,
					LookAt:   lookAt.Position(),
					LookFrom: lookFrom.Position(),
				}

				return nil, nil
			},
		},
	}

	allBuiltins := append(tengo.GetAllBuiltinFunctions(), customBuiltins...)
	allBuiltins = append(allBuiltins, builtins.SceneEngineBuiltins...)
	for i, b := range allBuiltins {
		symbolTable.DefineBuiltin(i, b.Name)
	}

	source, err := os.ReadFile(sceneScript)
	if err != nil {
		return err
	}

	moduleMap := tengo.NewModuleMap()
	addStdLib(moduleMap, "fmt", "math", "rand")
	moduleMap.AddSourceModule("userscript", source)
	//libraries.AddSceneEngineLibraries(moduleMap)

	fileSet := parser.NewFileSet()
	runtimeFile := fileSet.AddFile("runtime", -1, len(runtimeSource))
	p := parser.NewParser(runtimeFile, runtimeSource, nil)
	file, err := p.ParseFile()
	if err != nil {
		return err
	}

	c := tengo.NewCompiler(runtimeFile, symbolTable, nil, moduleMap, nil)
	if err := c.Compile(file); err != nil {
		return err
	}

	bytecode := c.Bytecode()
	machine := tengo.NewVM(bytecode, allBuiltins, globals, -1)
	if err := machine.Run(); err != nil {
		return err
	}

	return nil
}

func addStdLib(m *tengo.ModuleMap, libNames ...string) {
	for _, lib := range libNames {
		m.AddBuiltinModule(lib, stdlib.BuiltinModules[lib])
	}
}

func getConfig(o *tengo.Map) (messages.Config, error) {
	var config messages.Config
	config.FrameCount = 1
	config.Depth = 50
	config.Samples = 50
	config.UseBvh = true

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
		case "use_bvh":
			useBVH, ok := tengo.ToBool(val)
			if !ok {
				return config, fmt.Errorf("%w: got %T for %s", ErrConfigValueIncorrectType, val, key)
			}
			config.UseBvh = useBVH
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
