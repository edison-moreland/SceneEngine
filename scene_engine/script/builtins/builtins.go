package builtins

import (
	"github.com/d5/tengo/v2"
)

func GetArg[T any](toT func(object tengo.Object) (T, bool), args []tengo.Object, idx int, name string) (T, error) {
	arg, ok := toT(args[idx])
	if !ok {
		return arg, tengo.ErrInvalidArgumentType{
			Name:  name,
			Found: args[idx].TypeName(),
		}
	}

	return arg, nil
}

var SceneEngineBuiltins = []*tengo.BuiltinFunction{
	// Vec3
	{Name: "vec3", Value: builtinVec3},
	// Color
	{Name: "color", Value: builtinColor},
	// Shapes
	{Name: "sphere", Value: builtinSphere},
	// Materials
	{Name: "diffuse", Value: builtinDiffuse},
	{Name: "metallic", Value: builtinMetallic},
	{Name: "dielectric", Value: builtinDielectric},
	{Name: "emissive", Value: builtinEmissive},
	// Textures
	{Name: "uniform", Value: builtinUniform},
	{Name: "checker", Value: builtinChecker},
	{Name: "perlin", Value: builtinPerlin},
	// Rand
	{Name: "rand_color", Value: builtinRandColor},
	{Name: "rand_vec3", Value: builtinRandVec3},
	{Name: "rand_float", Value: builtinRandFloat},
	{Name: "rand_choice", Value: builtinRandChoice},
}
