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
	"testing"

	"github.com/gdamore/tcell"
)

type CanvasTest struct {
	screen tcell.Screen
	brush  tcell.Style
}

func (c CanvasTest) SetBrush(b tcell.Style) {
	c.brush = b
}

func (c *CanvasTest) UpdateBounds(r Rect) {

}

func (c CanvasTest) CreateCanvasFrom(r Rect) TCanvas {
	return NewCanvas(&c, r)
}

func (c CanvasTest) PrintChar(x int, y int, char rune) {
	c.screen.SetContent(x, y, char, nil, c.brush)
}

func (c CanvasTest) PrintCharWithBrush(x int, y int, char rune, brush tcell.Style) {
	c.screen.SetContent(x, y, char, nil, brush)
}

func (c CanvasTest) Fill(bounds Rect) {
	for y := bounds.Y; y < bounds.Height; y++ {
		for x := bounds.X; x < bounds.Width; x++ {
			c.PrintChar(x, y, ' ')
		}
	}
}

func TestCanvas_Create_sub_canvas(t *testing.T) {
	rootCanvas := CanvasTest{
		screen: tcell.NewSimulationScreen(""),
		brush: tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorBlack),
	}

	if e := rootCanvas.screen.Init(); e != nil {
		t.Errorf("Can't init screen: %+v\n", e)
	}

	rootCanvas.screen.Clear()

	childCanvas1DrawBounds := Rect{
		X:      5,
		Y:      5,
		Width:  40,
		Height: 20,
	}

	childCanvas1 := rootCanvas.CreateCanvasFrom(childCanvas1DrawBounds)

	if childCanvas1.(*Canvas).offset != childCanvas1DrawBounds {
		t.Errorf("Child Canvas 1 must be %+v and there is %+v", childCanvas1DrawBounds, childCanvas1.(*Canvas).offset)
	}

	childCanvas2Bounds := Rect{
		X:      3,
		Y:      2,
		Width:  10,
		Height: 2,
	}

	childCanvas2DrawBounds := Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 2,
	}

	childCanvas2 := childCanvas1.CreateCanvasFrom(childCanvas2Bounds)

	if childCanvas2.(*Canvas).offset != childCanvas2Bounds {
		t.Errorf("Child Canvas 1 must be %+v and there is %+v", childCanvas2Bounds, childCanvas2.(*Canvas).offset)
	}

	if childCanvas2.(*Canvas).draw != childCanvas2DrawBounds {
		t.Errorf("Child Canvas 1 must be %+v and there is %+v", childCanvas2DrawBounds, childCanvas2.(*Canvas).draw)
	}

	st := tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlue)

	childCanvas2.SetBrush(st)

	childCanvas2.PrintChar(1, 1, '&')

	if e := checkCell(rootCanvas.screen, 9, 8, '&', st, t); e != nil {
		t.Error(e)
	}
}
