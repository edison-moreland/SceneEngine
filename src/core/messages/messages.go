// Code generated by submsg; DO NOT EDIT.
package messages

import (
	"github.com/edison-moreland/SceneEngine/submsg/runtime/go"
	v5 "github.com/vmihailenco/msgpack/v5"
)

const (
	EngineMsgCoreReady  submsg.MsgId = 0
	EngineMsgCoreInfo   submsg.MsgId = 1
	EngineMsgPixelBatch submsg.MsgId = 2
)

type EngineServer interface {
	CoreReady(body []byte) error
	CoreInfo(body []byte) error
	PixelBatch(body []byte) error
}

func EngineRouter(s EngineServer) submsg.MsgReceiver {
	return func(id submsg.MsgId, body []byte) error {
		switch id {
		case EngineMsgCoreReady:
			return s.CoreReady(body)
		case EngineMsgCoreInfo:
			return s.CoreInfo(body)
		case EngineMsgPixelBatch:
			return s.PixelBatch(body)
		default:
			return submsg.ErrMsgIdUnknown
		}
	}
}

type EngineClient struct {
	s submsg.MsgSender
}

func NewEngineClient(s submsg.MsgSender) *EngineClient {
	return &EngineClient{s: s}
}
func (c *EngineClient) CoreReady(b []byte) {
	c.s(EngineMsgCoreReady, b)
}
func (c *EngineClient) CoreInfo(b []byte) {
	c.s(EngineMsgCoreInfo, b)
}
func (c *EngineClient) PixelBatch(b []byte) {
	c.s(EngineMsgPixelBatch, b)
}

const (
	CoreMsgInfo        submsg.MsgId = 0
	CoreMsgConfig      submsg.MsgId = 1
	CoreMsgRenderFrame submsg.MsgId = 2
)

type CoreServer interface {
	Info(body []byte) error
	Config(body []byte) error
	RenderFrame(body []byte) error
}

func CoreRouter(s CoreServer) submsg.MsgReceiver {
	return func(id submsg.MsgId, body []byte) error {
		switch id {
		case CoreMsgInfo:
			return s.Info(body)
		case CoreMsgConfig:
			return s.Config(body)
		case CoreMsgRenderFrame:
			return s.RenderFrame(body)
		default:
			return submsg.ErrMsgIdUnknown
		}
	}
}

type CoreClient struct {
	s submsg.MsgSender
}

func NewCoreClient(s submsg.MsgSender) *CoreClient {
	return &CoreClient{s: s}
}
func (c *CoreClient) Info(b []byte) {
	c.s(CoreMsgInfo, b)
}
func (c *CoreClient) Config(b []byte) {
	c.s(CoreMsgConfig, b)
}
func (c *CoreClient) RenderFrame(b []byte) {
	c.s(CoreMsgRenderFrame, b)
}

type MsgCoreInfo struct {
	Version string
}
type Config struct {
	AspectRatio float64
	Depth       uint64
	ImageHeight uint64
	ImageWidth  uint64
	Samples     uint64
}
type Position struct {
	X float64
	Y float64
	Z float64
}
type Color struct {
	B byte
	G byte
	R byte
}
type Pixel struct {
	Color Color
	X     uint64
	Y     uint64
}
type Lambert struct {
	Albedo Color
}
type Metal struct {
	Albedo  Color
	Scatter float64
}
type Dielectric struct {
	IndexOfRefraction float64
}
type Material struct {
	OneOf any
}

func MaterialFrom[T Lambert | Metal | Dielectric](v T) Material {
	return Material{OneOf: v}
}
func (o *Material) EncodeMsgpack(e *v5.Encoder) error {
	var err error
	switch o.OneOf.(type) {
	case Lambert:
		err = e.EncodeUint8(0)
	case Metal:
		err = e.EncodeUint8(1)
	case Dielectric:
		err = e.EncodeUint8(2)
	default:
		err = submsg.ErrUnknownOneOfField
	}
	if err != nil {
		return err
	}
	return e.Encode(o.OneOf)
}
func (o *Material) DecodeMsgpack(d *v5.Decoder) error {
	t, err := d.DecodeUint8()
	if err != nil {
		return err
	}
	switch t {
	case 0:
		var v Lambert
		err = d.Decode(&v)
		o.OneOf = v
	case 1:
		var v Metal
		err = d.Decode(&v)
		o.OneOf = v
	case 2:
		var v Dielectric
		err = d.Decode(&v)
		o.OneOf = v
	default:
		err = submsg.ErrUnknownOneOfField
	}
	return err
}

type Sphere struct {
	Origin Position
	Radius float64
}
type Shape struct {
	OneOf any
}

func ShapeFrom[T Sphere](v T) Shape {
	return Shape{OneOf: v}
}
func (o *Shape) EncodeMsgpack(e *v5.Encoder) error {
	var err error
	switch o.OneOf.(type) {
	case Sphere:
		err = e.EncodeUint8(0)
	default:
		err = submsg.ErrUnknownOneOfField
	}
	if err != nil {
		return err
	}
	return e.Encode(o.OneOf)
}
func (o *Shape) DecodeMsgpack(d *v5.Decoder) error {
	t, err := d.DecodeUint8()
	if err != nil {
		return err
	}
	switch t {
	case 0:
		var v Sphere
		err = d.Decode(&v)
		o.OneOf = v
	default:
		err = submsg.ErrUnknownOneOfField
	}
	return err
}

type Object struct {
	Material Material
	Shape    Shape
}
type Camera struct {
	Aperture float64
	Fov      float64
	LookAt   Position
	LookFrom Position
}
type Scene struct {
	Camera  Camera
	Objects []Object
}
