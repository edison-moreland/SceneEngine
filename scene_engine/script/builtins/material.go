package builtins

import (
	"fmt"

	"github.com/d5/tengo/v2"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func builtinDiffuse(args ...tengo.Object) (ret tengo.Object, err error) {
	// Diffuse takes one argument, a texture
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	texture, err := GetArg(ToTexture, args, 0, "texture")
	if err != nil {
		return nil, err
	}

	return &Material{Material: messages.MaterialFrom(messages.Diffuse{
		Texture: texture,
	})}, nil
}

func builtinMetallic(args ...tengo.Object) (ret tengo.Object, err error) {
	// Metal takes two arguments, a texture and a scatter
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	texture, err := GetArg(ToTexture, args, 0, "texture")
	if err != nil {
		return nil, err
	}

	scatter, err := GetArg(tengo.ToFloat64, args, 1, "scatter")
	if err != nil {
		return nil, err
	}

	return &Material{Material: messages.MaterialFrom(messages.Metallic{
		Texture: texture,
		Scatter: scatter,
	})}, nil
}

func builtinDielectric(args ...tengo.Object) (ret tengo.Object, err error) {
	// Dielectric takes one argument, the index of refraction
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	ior, err := GetArg(tengo.ToFloat64, args, 0, "ior")
	if err != nil {
		return nil, err
	}

	return &Material{Material: messages.MaterialFrom(messages.Dielectric{
		IndexOfRefraction: ior,
	})}, nil
}

func builtinEmissive(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	texture, err := GetArg(ToTexture, args, 0, "texture")
	if err != nil {
		return nil, err
	}

	brightness, err := GetArg(tengo.ToFloat64, args, 1, "brightness")
	if err != nil {
		return nil, err
	}

	return &Material{Material: messages.MaterialFrom(messages.Emissive{
		Texture:    texture,
		Brightness: brightness,
	})}, nil
}

func ToMaterial(o tengo.Object) (messages.Material, bool) {
	switch o := o.(type) {
	case *Material:
		return o.Material, true
	}

	return messages.Material{}, false
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
