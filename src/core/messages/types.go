package messages

type MsgCoreInfo struct {
	Version string
}

type Config struct {
	ImageWidth  uint64
	ImageHeight uint64
}

type Color struct {
	R float64
	G float64
	B float64
}

type Pixel struct {
	X     uint64
	Y     uint64
	Color Color
}
