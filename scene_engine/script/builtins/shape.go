package builtins

import (
	"fmt"

	"github.com/d5/tengo/v2"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func builtinSphere(args ...tengo.Object) (ret tengo.Object, err error) {
	// Sphere takes two arguments, origin and radius
	if len(args) != 2 {
		return nil, tengo.ErrWrongNumArguments
	}

	origin, err := GetArg(ToPosition, args, 0, "origin")
	if err != nil {
		return nil, err
	}

	radius, err := GetArg(tengo.ToFloat64, args, 1, "radius")
	if err != nil {
		return nil, err
	}

	return &Shape{Shape: messages.ShapeFrom(messages.Sphere{
		Origin: origin,
		Radius: radius,
	})}, nil
}

func ToShape(o tengo.Object) (messages.Shape, bool) {
	switch o := o.(type) {
	case *Shape:
		return o.Shape, true
	}

	return messages.Shape{}, false
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
