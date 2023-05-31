package libraries

import (
	"fmt"

	"github.com/d5/tengo/v2"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

var ShapeModule = map[string]tengo.Object{
	"Sphere": &tengo.UserFunction{Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
		// Sphere takes two arguments, origin and radius
		if len(args) != 2 {
			return nil, tengo.ErrWrongNumArguments
		}

		origin, ok := args[0].(*Vec3)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "origin"}
		}

		radius, ok := tengo.ToFloat64(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{Name: "radius"}
		}

		return &Shape{Shape: messages.ShapeFrom(messages.Sphere{
			Origin: origin.Position(),
			Radius: radius,
		})}, nil
	}},
}

type Shape struct {
	tengo.ObjectImpl
	messages.Shape
}

func (s *Shape) TypeName() string {
	return "Shape"
}

func (s *Shape) String() string {
	return fmt.Sprintf("Shape(%v)", s.OneOf)
}
