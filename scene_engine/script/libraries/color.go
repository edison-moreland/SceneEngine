package libraries

import (
	"fmt"

	"github.com/d5/tengo/v2"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

var ColorModule = map[string]tengo.Object{
	"New": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}

		r, ok := tengo.ToInt(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "R"}
		}

		g, ok := tengo.ToInt(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "G"}
		}

		b, ok := tengo.ToInt(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "B"}
		}

		return &Color{Value: messages.Color{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
		}}, nil
	}},

	"Random": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		saturation := float32(0.5)
		switch len(args) {
		case 1:
			newSaturation, ok := tengo.ToFloat64(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{Name: "saturation"}
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
	}},

	"Red": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 255, G: 0, B: 0}}, nil
	}},

	"Green": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 0, G: 255, B: 0}}, nil
	}},

	"Blue": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 0, G: 0, B: 255}}, nil
	}},

	"Yellow": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 255, G: 255, B: 0}}, nil
	}},

	"Cyan": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 0, G: 255, B: 255}}, nil
	}},

	"Magenta": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 255, G: 0, B: 255}}, nil
	}},

	"Black": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 0, G: 0, B: 0}}, nil
	}},

	"White": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		return &Color{Value: messages.Color{R: 255, G: 255, B: 255}}, nil
	}},
}

type Color struct {
	tengo.ObjectImpl
	Value messages.Color
}

func (c *Color) TypeName() string {
	return "Color"
}

func (c *Color) String() string {
	return fmt.Sprintf("Color(%d, %d, %d", c.Value.R, c.Value.G, c.Value.B)
}
