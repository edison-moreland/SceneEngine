package builtins

import (
	"github.com/d5/tengo/v2"
)

var SceneEngineBuiltins = []*tengo.BuiltinFunction{
	// Vec3
	{Name: "vec3", Value: builtinVec3},
	// Color
	{Name: "color", Value: builtinColor},
	// Shapes
	{Name: "sphere", Value: builtinSphere},
	// Materials
	{Name: "lambert", Value: builtinLambert},
	{Name: "metal", Value: builtinMetal},
	{Name: "dielectric", Value: builtinDielectric},
	// Rand
	{Name: "rand_color", Value: builtinRandColor},
	{Name: "rand_vec3", Value: builtinRandVec3},
	{Name: "rand_float", Value: builtinRandFloat},
	{Name: "rand_choice", Value: builtinRandChoice},
}
