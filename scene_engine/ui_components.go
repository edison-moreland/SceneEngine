package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawInfo(row int32, key, value string) {
	textSize := int32(10)
	margin := int32(2)

	rl.DrawText(
		fmt.Sprintf("%s: %s", key, value),
		margin, margin+(textSize*row), textSize, rl.Black,
	)
}

type Button struct {
	rec rl.Rectangle

	text  string
	color rl.Color

	textSize   int32
	textMargin int32
}

func NewButton(text string, color rl.Color, x, y float32) Button {
	b := Button{
		text:       text,
		color:      color,
		textSize:   20,
		textMargin: 5,
	}

	textWidth := rl.MeasureText(text, b.textSize)

	b.rec = rl.NewRectangle(
		x, y,
		float32(textWidth+(b.textMargin*2)),
		float32(b.textSize+(b.textMargin*2)),
	)

	return b
}

func (b *Button) Draw() {
	rl.DrawRectangleRec(b.rec, rl.Red)
	rl.DrawText(
		b.text,
		int32(b.rec.X)+b.textMargin,
		int32(b.rec.Y)+b.textMargin,
		b.textSize,
		rl.Black,
	)
}

func (b *Button) Down() bool {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()
		return rl.CheckCollisionPointRec(mousePos, b.rec)
	}

	return false
}

type Slider struct {
	backRec, handleRec     rl.Rectangle
	handleStart, handleEnd rl.Vector2

	margin float32

	valueMax, valueMin float64
	value              float64
}

func NewSlider(x, y, width, height float32, min, max float64) Slider {
	s := Slider{
		backRec:  rl.NewRectangle(x, y, width, height),
		margin:   5,
		valueMin: min,
		valueMax: max,
		value:    min,
	}

	s.handleRec = rl.NewRectangle(
		x+s.margin, y+s.margin,
		5, height-(s.margin*2),
	)

	s.handleStart = rl.NewVector2(
		s.handleRec.X+(s.handleRec.Width/2),
		s.handleRec.Y+(s.handleRec.Height/2),
	)

	s.handleEnd = rl.NewVector2(
		x+(width-((s.handleStart.X-x)*2)),
		s.handleStart.Y,
	)

	return s
}

func (s *Slider) Draw() {
	rl.DrawRectangleRec(s.backRec, rl.Gray)
	rl.DrawLineV(s.handleStart, s.handleEnd, rl.DarkGray)
	rl.DrawRectangleRec(s.handleRec, rl.Red)
}
