package builtins

import (
	"github.com/d5/tengo/v2"
)

var SceneEngineBuiltins = []*tengo.BuiltinFunction{
	// Vec3
	{Name: "vec3", Value: builtinVec3},
	// Color
	{Name: "color", Value: builtinColor},
	{Name: "color_rand", Value: builtinRandColor},
	// Shapes
	{Name: "sphere", Value: builtinSphere},
	// Materials
	{Name: "lambert", Value: builtinLambert},
	{Name: "metal", Value: builtinMetal},
	{Name: "dielectric", Value: builtinDielectric},
}
