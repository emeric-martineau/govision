package base

// Copyright 2020 The GoVision Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"github.com/gdamore/tcell"
)

// Canvas is virtal screen to draw components.
type Canvas struct {
	// Parent canvas.
	parent TCanvas
	// Brush styte
	brush tcell.Style
	// Offset to convert local position to parent position.
	offset Rect
	// Draw area.
	draw Rect
}

// SetBrush set the brush to draw.
func (c *Canvas) SetBrush(b tcell.Style) {
	c.brush = b
}

// UpdateBounds call when component move or resize.
func (c *Canvas) UpdateBounds(r Rect) {
	c.offset = r

	c.draw = Rect{
		X:      0,
		Y:      0,
		Width:  r.Width,
		Height: r.Height,
	}
}

// CreateCanvasFrom create a sub-canvas for `r` parameter.
func (c *Canvas) CreateCanvasFrom(r Rect) TCanvas {
	return &Canvas{
		parent: c,
		offset: r,
		draw: Rect{
			X:      0,
			Y:      0,
			Width:  r.Width,
			Height: r.Height,
		},
	}
}

// PrintChar print a charactere.
func (c *Canvas) PrintChar(x int, y int, char rune) {
	// Check if out of me.
	if InHorizontal(x, c.draw) && InVertical(y, c.draw) {
		c.parent.PrintCharWithBrush(x+c.offset.X, y+c.offset.Y, char, c.brush)
	}
}

// PrintCharWithBrush print a charactere with brush.
func (c *Canvas) PrintCharWithBrush(x int, y int, char rune, brush tcell.Style) {
	// Check if out of me.
	if InHorizontal(x, c.draw) && InVertical(y, c.draw) {
		c.parent.PrintCharWithBrush(x+c.offset.X, y+c.offset.Y, char, brush)
	}
}

// Fill zone of canvas.
func (c *Canvas) Fill(bounds Rect) {
	for y := bounds.Y; y < bounds.Height; y++ {
		for x := bounds.X; x < bounds.Width; x++ {
			c.PrintChar(x, y, ' ')
		}
	}
}

//------------------------------------------------------------------------------
// Constrcutor.

// NewCanvas create a sub canvas.
func NewCanvas(parent TCanvas, r Rect) TCanvas {
	return &Canvas{
		parent: parent,
		offset: r,
		draw: Rect{
			X:      0,
			Y:      0,
			Width:  r.Width,
			Height: r.Height,
		},
	}
}
