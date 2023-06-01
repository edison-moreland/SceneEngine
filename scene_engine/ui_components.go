package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/edison-moreland/SceneEngine/scene_engine/layout"
)

type TextBlock struct {
	lines    []string
	textSize int32

	p *layout.Rec
	d layout.Direction
}

func NewTextBlock(lines []string, textSize int32, p *layout.Rec, d layout.Direction) TextBlock {
	return TextBlock{lines: lines, textSize: textSize, p: p, d: d}
}

func (t *TextBlock) Measure() *layout.Rec {
	bounds := rl.Rectangle{}
	bounds.Height = float32(len(t.lines) * int(t.textSize))

	for _, line := range t.lines {
		width := float32(rl.MeasureText(line, t.textSize))
		if width > bounds.Width {
			bounds.Width = width
		}
	}

	return layout.Layout(bounds)
}

func (t *TextBlock) Draw() {
	rec := t.Measure().Layout(t.p, t.d)

	for row, line := range t.lines {
		rl.DrawText(
			line,
			int32(rec.X), int32(rec.Y)+(t.textSize*int32(row)), t.textSize, rl.Black,
		)
	}
}

type StatBlock struct {
	TextBlock
}

func NewStatBlock(p *layout.Rec, d layout.Direction) StatBlock {
	return StatBlock{
		TextBlock{textSize: 10, p: p, d: d},
	}
}

func (s *StatBlock) AddStat(name, value string) {
	s.lines = append(s.lines, fmt.Sprintf("%s: %s", name, value))
}

type Button struct {
	rec rl.Rectangle

	text  string
	color rl.Color

	textSize   int32
	textMargin int32
}

func NewButton(text string, color rl.Color, p *layout.Rec, d layout.Direction) Button {
	b := Button{
		text:       text,
		color:      color,
		textSize:   20,
		textMargin: 5,
	}

	textWidth := rl.MeasureText(text, b.textSize)

	b.rec = rl.NewRectangle(
		0, 0,
		float32(textWidth+(b.textMargin*2)),
		float32(b.textSize+(b.textMargin*2)),
	)

	// TODO: Rename some functions?
	b.rec = layout.Layout(b.rec).Layout(p, d).Rectangle

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

func NewSlider(rec *layout.Rec, min, max float64) Slider {
	s := Slider{
		backRec:  rec.Rectangle,
		margin:   5,
		valueMin: min,
		valueMax: max,
		value:    min,
	}

	handleWidth := float32(5)
	s.handleRec = rec.Margin(s.margin).Resize(layout.Horizontal, layout.West, handleWidth).Rectangle

	s.handleStart = layout.CenterPoint(s.handleRec)

	timelineWidth := s.backRec.Width - ((s.handleStart.X - s.backRec.X) * 2)
	s.handleEnd = s.handleStart
	s.handleEnd.X += timelineWidth

	return s
}

func (s *Slider) Draw() {
	rl.DrawRectangleRec(s.backRec, rl.Gray)
	rl.DrawLineV(s.handleStart, s.handleEnd, rl.DarkGray)
	rl.DrawRectangleRec(s.handleRec, rl.Red)
}
