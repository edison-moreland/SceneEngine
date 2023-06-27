package builtins

import (
	"fmt"
	"math/rand"

	"github.com/d5/tengo/v2"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

// TODO: Make rand deterministic per scene?

func builtinRandFloat(args ...tengo.Object) (ret tengo.Object, err error) {
	return &tengo.Float{Value: rand.Float64()}, nil
}

func builtinRandVec3(args ...tengo.Object) (ret tengo.Object, err error) {
	return newVec3(rl.NewVector3(
		rand.Float32(),
		rand.Float32(),
		rand.Float32(),
	)), nil
}

func builtinRandColor(args ...tengo.Object) (ret tengo.Object, err error) {
	saturation := float32(0.5)
	switch len(args) {
	case 1:

		newSaturation, err := GetArg(tengo.ToFloat64, args, 0, "saturation")
		if err != nil {
			return nil, err
		}
		saturation = float32(newSaturation)
	case 0:
		break
	default:
		return nil, tengo.ErrWrongNumArguments
	}

	hue := rl.GetRandomValue(0, 360)

	color := rl.ColorFromHSV(float32(hue), saturation, 1.0)

	return &Color{Value: messages.Color{
		R: color.R,
		G: color.G,
		B: color.B,
	}}, nil
}

func builtinRandChoice(args ...tengo.Object) (ret tengo.Object, err error) {
	// rand_choice(
	//   [0.80, lambert(rand_color())],
	//   [0.15, metal(rand_color(0.5), rand_float()/2)],
	//   [0.05, dielectric(1.5)]
	// )
	if len(args) < 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	totalWeight := float64(0)
	weights := make([]float64, len(args))
	choices := make([]tengo.Object, len(args))

	for i, arg := range args {
		argArr, ok := arg.(*tengo.Array)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: fmt.Sprintf("args[%d]", i)}
		}

		if len(argArr.Value) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}

		weight, err := GetArg(tengo.ToFloat64, argArr.Value, 0, fmt.Sprintf("args[%d].weight", i))
		if err != nil {
			return nil, err
		}

		totalWeight += weight
		weights[i] = weight
		choices[i] = argArr.Value[1]
	}
	if totalWeight != 1.0 {
		return nil, fmt.Errorf("weights must add up to 1")
	}

	random := rand.Float64()
	count := float64(0)
	for i, weight := range weights {
		count += weight

		if random < count {
			return choices[i], nil
		}
	}

	panic("This should not happen :)")
}
