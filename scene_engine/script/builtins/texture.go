package builtins

import (
	"fmt"

	"github.com/d5/tengo/v2"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func builtinUniform(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 1 {
		return nil, tengo.ErrWrongNumArguments
	}

	color, err := GetArg(ToColor, args, 0, "color")
	if err != nil {
		return nil, err
	}

	return &Texture{Texture: messages.TextureFrom(messages.Uniform{
		Color: color,
	})}, nil
}

func builtinChecker(args ...tengo.Object) (tengo.Object, error) {
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	even, err := GetArg(ToTexture, args, 0, "even")
	if err != nil {
		return nil, err
	}

	odd, err := GetArg(ToTexture, args, 1, "odd")
	if err != nil {
		return nil, err
	}

	return &Texture{Texture: messages.TextureFrom(messages.Checker{
		Even: even,
		Odd:  odd,
	})}, nil
}

func ToTexture(o tengo.Object) (messages.Texture, bool) {
	switch o := o.(type) {
	case *Texture:
		return o.Texture, true
	case *Color:
		return messages.TextureFrom(messages.Uniform{Color: o.Value}), true
	}

	return messages.Texture{}, false
}

type Texture struct {
	tengo.ObjectImpl
	messages.Texture
}

func (t *Texture) TypeName() string {
	return "Texture"
}

func (t *Texture) String() string {
	return fmt.Sprintf("Texture(%v)", t.OneOf)
}
