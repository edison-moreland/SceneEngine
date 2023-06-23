package builtins

import (
	"fmt"

	"github.com/d5/tengo/v2"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func builtinLambert(args ...tengo.Object) (ret tengo.Object, err error) {
	// Lambert takes one argument, a color
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	albedo, ok := args[0].(*Color)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "albedo",
			Expected: "Color",
			Found:    args[0].TypeName(),
		}
	}

	return &Material{Material: messages.MaterialFrom(messages.Lambert{
		Albedo: albedo.Value,
	})}, nil
}

func builtinMetal(args ...tengo.Object) (ret tengo.Object, err error) {
	// Metal takes two arguments, a color and a scatter
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	albedo, ok := args[0].(*Color)
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "albedo",
			Expected: "Color",
			Found:    args[0].TypeName(),
		}
	}

	scatter, ok := tengo.ToFloat64(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "scatter",
			Expected: "Float",
			Found:    args[0].TypeName(),
		}
	}

	return &Material{Material: messages.MaterialFrom(messages.Metal{
		Albedo:  albedo.Value,
		Scatter: scatter,
	})}, nil
}

func builtinDielectric(args ...tengo.Object) (ret tengo.Object, err error) {
	// Dielectric takes one argument, the index of refraction
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	ior, ok := tengo.ToFloat64(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "ior",
			Expected: "Float",
			Found:    args[0].TypeName(),
		}
	}

	return &Material{Material: messages.MaterialFrom(messages.Dielectric{
		IndexOfRefraction: ior,
	})}, nil
}
	
type Material struct {
	tengo.ObjectImpl
	messages.Material
}

func (m *Material) TypeName() string {
	return "Material"
}

func (m *Material) String() string {
	return fmt.Sprintf("Material(%v)", m.OneOf)
}
