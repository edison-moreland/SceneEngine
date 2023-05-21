package messages

type MsgCoreInfo struct {
	Version string
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
