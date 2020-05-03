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

// PrintStringOnScreen helper to print string. For Debug only.
func PrintStringOnScreen(s tcell.Screen, bc tcell.Color, fc tcell.Color, x int, y int, msg string) {
	style := tcell.StyleDefault.
		Foreground(fc).
		Background(bc)

	for pos, char := range msg {
		s.
			SetContent(x+pos, y, char, nil, style)
	}
}

// MaxInt return greater integer.
func MaxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt return less integer.
func MinInt(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func inHorizontal(x int, r Rect) bool {
	return (x >= r.X) && (x < r.X+r.Width)
}

func inVertical(y int, r Rect) bool {
	return (y >= r.Y) && (y < r.Y+r.Height)
}

// Intersect return interesection of rectangle.
// If no intersection, return Rect{X:0, Y:0, Width:0, Height: 0}.
func Intersect(r1 Rect, r2 Rect) Rect {
	var leftX int
	var width int
	//rightX := MinInt(r1.X+r1.Width, r2.X+r2.Width)
	var height int
	var topY int

	if inHorizontal(r2.X, r1) || inHorizontal(r1.X, r2) {
		leftX = MaxInt(r1.X, r2.X)
		rightX := MinInt(r1.X+r1.Width, r2.X+r2.Width)
		width = MaxInt(0, rightX-leftX)
	} else {
		leftX = 0
		width = 0
	}

	if inVertical(r2.Y, r1) || inVertical(r1.Y, r2) {
		topY = MaxInt(r1.Y, r2.Y)
		bottomY := MinInt(r1.Y+r1.Height, r2.Y+r2.Height)
		height = MaxInt(0, bottomY-topY)
	} else {
		topY = 0
		height = 0
	}

	return Rect{
		X:      leftX,
		Y:      topY,
		Width:  width,
		Height: height,
	}
}

// CalculateAbsolutePosition helper to calculate position of component from
// relative position of parent, to absolute position of screen.
// If component is not visible, return Rect{X:0, Y:0, Width:0, Height: 0}.
func CalculateAbsolutePosition(view TView) Rect {
	viewRect := view.GetBounds()

	p := view.GetParent()

	var v TView

	switch p := p.(type) {
	case TView:
		v = p

	default:
		v = nil
	}

	for v != nil {
		bounds := v.GetBounds()
		clientBounds := v.GetClientBounds()

		// viewRect must be have same referencial as parent
		viewRect = Rect{
			X:      viewRect.X + bounds.X + clientBounds.X,
			Y:      viewRect.Y + bounds.Y + clientBounds.Y,
			Width:  viewRect.Width,
			Height: viewRect.Height,
		}

		// parentBound must take care of client bound
		parentBound := Rect{
			X:      bounds.X + clientBounds.X,
			Y:      bounds.Y + clientBounds.Y,
			Width:  MinInt(bounds.Width, clientBounds.Width),
			Height: MinInt(bounds.Height, clientBounds.Height),
		}

		viewRect = Intersect(viewRect, parentBound)

		if viewRect.Width == 0 && viewRect.Height == 0 {
			// Component is not visible, stop now.
			return viewRect
		}

		p := v.GetParent()

		switch p := p.(type) {
		case TView:
			v = p
		default:
			v = nil
		}
	}

	return viewRect
}
