package layout

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Axis int

const (
	Horizontal Axis = iota
	Vertical
)

type Direction int

const (
	North Direction = iota
	South
	East
	West
	NorthWest
	NorthEast
	SouthWest
	SouthEast
)

// These primitives are used to layout a rectangle within another rectangle

func CenterPoint(r rl.Rectangle) rl.Vector2 {
	return rl.NewVector2(
		r.X+(r.Width/2),
		r.Y+(r.Height/2),
	)
}

type Rec struct {
	rl.Rectangle
}

func Layout(r rl.Rectangle) *Rec {
	return &Rec{
		Rectangle: r,
	}
}

// Window creates a layout rec of the current window's size
func Window() *Rec {
	return &Rec{
		Rectangle: rl.NewRectangle(
			0, 0,
			float32(rl.GetScreenWidth()),
			float32(rl.GetScreenHeight()),
		),
	}
}

func (l *Rec) Margin(d float32) *Rec {
	// Inset the entire rec by d
	return &Rec{
		Rectangle: rl.NewRectangle(
			l.X+d, l.Y+d,
			l.Width-(d*2),
			l.Height-(d*2),
		),
	}
}

func (l *Rec) Resize(axis Axis, direction Direction, value float32) *Rec {
	rec := l.Rectangle

	switch axis {
	case Vertical:
		rec.Height = value
		difference := l.Height - value

		switch direction {
		case North:
			// Nothing to do

		case South:
			rec.Y += difference

		default:
			panic("Incompatible direction with axis")
		}

	case Horizontal:
		rec.Width = value
		difference := l.Width - value

		switch direction {
		case West:
			// Nothing to do

		case East:
			rec.X += difference

		default:
			panic("Incompatible direction with axis")
		}

	}

	return &Rec{
		Rectangle: rec,
	}
}

func (l *Rec) Layout(p *Rec, d Direction) *Rec {
	if (l.Width > p.Width) || (l.Height > p.Height) {
		panic("lrec must fit within prec")
	}

	// Start with rec l floating in center of rec p
	// Move rec l in direction d until it hits an edge of rec p

	rec := l.Rectangle

	switch d {
	case NorthWest:
		rec.X = p.X
		rec.Y = p.Y
	case NorthEast:
		rec.X = p.X - rec.Width
		rec.Y = p.Y
	case SouthWest:
		rec.X = p.X
		rec.Y = p.Y - rec.Height
	case SouthEast:
		rec.X = p.X - rec.Width
		rec.Y = p.Y - rec.Height
	}

	return Layout(rec)
}
