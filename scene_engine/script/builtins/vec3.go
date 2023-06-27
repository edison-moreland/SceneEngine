package builtins

import (
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func builtinVec3(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 3 {
		return nil, tengo.ErrWrongNumArguments
	}

	x, err := getArg(tengo.ToFloat64, args, 0, "X")
	if err != nil {
		return nil, err
	}

	y, err := getArg(tengo.ToFloat64, args, 1, "Y")
	if err != nil {
		return nil, err
	}

	z, err := getArg(tengo.ToFloat64, args, 2, "Z")
	if err != nil {
		return nil, err
	}

	return newVec3(rl.NewVector3(
		float32(x),
		float32(y),
		float32(z),
	)), nil
}

func ToVector3(o tengo.Object) (rl.Vector3, bool) {
	switch o := o.(type) {
	case *Vec3:
		return o.Value, true
	}

	return rl.Vector3{}, false
}

func ToPosition(o tengo.Object) (messages.Position, bool) {
	switch o := o.(type) {
	case *Vec3:
		return o.Position(), true
	}

	return messages.Position{}, false
}

// Vec3 is a tengo vector
type Vec3 struct {
	tengo.ObjectImpl
	Value rl.Vector3
}

func newVec3(vector rl.Vector3) *Vec3 {
	return &Vec3{
		Value: vector,
	}
}

func (v *Vec3) TypeName() string {
	return "Vec3"
}

func (v *Vec3) String() string {
	return fmt.Sprintf("Vec3(%f, %f, %f)", v.Value.X, v.Value.Y, v.Value.Z)
}

func (v *Vec3) Copy() tengo.Object {
	return &Vec3{Value: rl.Vector3{
		X: v.Value.X,
		Y: v.Value.Y,
		Z: v.Value.Z,
	}}
}

func (v *Vec3) IsFalsy() bool {
	return rl.Vector3Length(v.Value) == 0
}

func (v *Vec3) Equals(rhs tengo.Object) bool {
	rhsVec, ok := rhs.(*Vec3)
	if !ok {
		return false
	}

	return (v.Value.X == rhsVec.Value.X) &&
		(v.Value.Y == rhsVec.Value.Y) &&
		(v.Value.Z == rhsVec.Value.Z)
}

func (v *Vec3) IndexGet(index tengo.Object) (tengo.Object, error) {
	idx, ok := tengo.ToString(index)
	if !ok {
		return nil, tengo.ErrInvalidIndexType
	}

	switch idx {
	case "x":
		return &tengo.Float{Value: float64(v.Value.X)}, nil
	case "y":
		return &tengo.Float{Value: float64(v.Value.Y)}, nil
	case "z":
		return &tengo.Float{Value: float64(v.Value.Z)}, nil
	case "length":
		return &tengo.Float{Value: float64(rl.Vector3Length(v.Value))}, nil
	default:
		return tengo.UndefinedValue, nil
	}
}

func (v *Vec3) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	rhsFloat, ok := tengo.ToFloat64(rhs)
	if ok {
		return v.floatOp(op, float32(rhsFloat))
	}

	rhsVec3, ok := rhs.(*Vec3)
	if ok {
		return v.vec3Op(op, rhsVec3.Value)
	}

	return nil, tengo.ErrNotImplemented
}

func (v *Vec3) floatOp(op token.Token, rhs float32) (tengo.Object, error) {
	switch op {
	case token.Mul:
		return newVec3(rl.Vector3Multiply(v.Value, rhs)), nil
	default:
		return nil, tengo.ErrInvalidOperator
	}
}

func (v *Vec3) vec3Op(op token.Token, rhs rl.Vector3) (tengo.Object, error) {
	switch op {
	case token.Add:
		return newVec3(rl.Vector3Add(v.Value, rhs)), nil
	case token.Sub:
		return newVec3(rl.Vector3Subtract(v.Value, rhs)), nil
	case token.Mul:
		return newVec3(rl.Vector3MultiplyV(v.Value, rhs)), nil
	default:
		return nil, tengo.ErrInvalidOperator
	}
}

func (v *Vec3) Position() messages.Position {
	return messages.Position{
		X: float64(v.Value.X),
		Y: float64(v.Value.Y),
		Z: float64(v.Value.Z),
	}
}
