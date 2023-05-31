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
