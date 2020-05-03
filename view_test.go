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

func TestView_Dummy_for_code_coverage(t *testing.T) {
	s := mkTestScreen(t, "")
	c := NewView("You know my name", s)
	c.SetFocused(true)
	c.GetFocused()
	c.SetVisible(true)
	c.GetVisible()
	c.SetBackgroundColor(tcell.ColorBlack)
	c.GetBackgroundColor()
	c.SetForegroundColor(tcell.ColorBlack)
	c.GetForegroundColor()
	c.GetScreen()
}

func TestView_draw_one_view(t *testing.T) {
	screenStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	s := mkTestScreen(t, "")
	defer s.Fini()

	s.SetStyle(screenStyle)

	s.Clear()

	st := tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlue)

	c := NewView("You know my name", s)
	c.SetBounds(Rect{
		X:      1,
		Y:      1,
		Width:  3,
		Height: 4,
	})
	c.SetForegroundColor(tcell.ColorYellow)
	c.SetBackgroundColor(tcell.ColorBlue)
	c.SetVisible(true)

	c.HandleMessage(Message{
		Handler: BroadcastHandler(),
		Type:    WmDraw,
	})

	s.Show()

	for x := 1; x < 4; x++ {
		for y := 1; y < 4; y++ {
			if e := checkCell(x, y, ' ', s, st, t); e != nil {
				t.Error(e)
			}
		}
	}
}

func TestView_draw_one_view_with_OnDraw(t *testing.T) {
	isCalled := false
	screenStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	s := mkTestScreen(t, "")
	defer s.Fini()

	s.SetStyle(screenStyle)

	if e := s.Init(); e != nil {
		t.Error("Can't init screen")
		return
	}

	s.Clear()

	c := NewView("You know my name", s)
	c.OnDraw = func(v TView, s tcell.Screen) {
		isCalled = true
	}

	c.SetBounds(Rect{
		X:      1,
		Y:      1,
		Width:  3,
		Height: 4,
	})
	c.SetForegroundColor(tcell.ColorYellow)
	c.SetBackgroundColor(tcell.ColorBlue)
	c.SetVisible(true)

	c.HandleMessage(Message{
		Handler: BroadcastHandler(),
		Type:    WmDraw,
	})

	s.Show()

	// Component not draw, cause OnDraw do nothing in our case
	for x := 1; x < 4; x++ {
		for y := 1; y < 4; y++ {
			if e := checkCell(x, y, ' ', s, 0, t); e != nil {
				t.Error(e)
			}
		}
	}

	if !isCalled {
		t.Error("OnDraw not have be called!")
	}
}

// Component 1 (c1):
//   X: 1
//   Y: 1
//   Width: 10
//   Height: 20
//
// Component 2 (c2):
//   X: 2
//   Y: 2
//   Width: 5
//   Height: 5
//
//
//  x:1 / y:1                                          x:11 / y:1
//  +-----------------------C1-------------------------+
//  |                                                  |
//  |  x:3 / y:3                          x:8 / y:3    |
//  |  +----------------C2----------------+            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  |                                  |            |
//  |  +----------------------------------+            |
//  |  x:3 / y:8                          x:8 / y:8    |
//  |                                                  |
//  |                                                  |
//  |                                                  |
//  +--------------------------------------------------+
//  x:1 / y: 21                                        x:11 / y: 21
func TestView_draw_one_view_in_another_view(t *testing.T) {
	screenStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	s := mkTestScreen(t, "")
	defer s.Fini()

	s.SetStyle(screenStyle)

	if e := s.Init(); e != nil {
		t.Error("Can't init screen")
		return
	}

	s.Clear()

	c1Style := tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlue)

	c2Style := tcell.StyleDefault.
		Foreground(tcell.ColorBlue).
		Background(tcell.ColorGreen)

	c1 := NewView("MainComponent", s)
	c1.SetBounds(Rect{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 20,
	})
	c1.SetForegroundColor(tcell.ColorYellow)
	c1.SetBackgroundColor(tcell.ColorBlue)
	c1.SetVisible(true)

	c2 := NewView("ChildComponent", s)
	c2.SetBounds(Rect{
		X:      2,
		Y:      2,
		Width:  5,
		Height: 5,
	})
	c2.SetForegroundColor(tcell.ColorBlue)
	c2.SetBackgroundColor(tcell.ColorGreen)
	c2.SetParent(&c1)
	c2.SetVisible(true)

	c1.AddChild(&c2)

	c1.HandleMessage(Message{
		Handler: BroadcastHandler(),
		Type:    WmDraw,
	})

	s.Show()

	// Check component 1 ---------------------------------------------------------
	// Top
	for x := 1; x < 11; x++ {
		for y := 1; y < 3; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}

	// Bottom
	for x := 1; x < 11; x++ {
		for y := 9; y < 21; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}

	// Left
	for x := 1; x < 2; x++ {
		for y := 1; y < 21; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}

	// Right
	for x := 9; x < 11; x++ {
		for y := 1; y < 21; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}

	// Check component 2 ---------------------------------------------------------
	for x := 3; x < 8; x++ {
		for y := 3; y < 8; y++ {
			if e := checkCell(x, y, ' ', s, c2Style, t); e != nil {
				t.Error(e)
			}
		}
	}
}

// Component 1 (c1):
//   X: 1
//   Y: 1
//   Width: 10
//   Height: 20
//
// Component 2 (c2):
//   X: -2
//   Y: -2
//   Width: 5
//   Height: 5
//
//
//  x:1 / y:1                                          x:11 / y:1
//  +-----------------------C1-------------------------+
//  |                         |                        |
//  |                         |                        |
//  |                         |                        |
//  |                         |                        |
//  |                         |                        |
//  |                         |                        |
//  |                         |                        |
//  |                         |                        |
//  |----------C2-------------+                        |
//  |                         x:4 / y:4                |
//  |                                                  |
//  |                                                  |
//  |                                                  |
//  +--------------------------------------------------+
//  x:1 / y: 21                                        x:11 / y: 21
func TestView_draw_one_view_in_another_view_partial_out(t *testing.T) {
	screenStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	s := mkTestScreen(t, "")
	defer s.Fini()

	s.SetStyle(screenStyle)

	if e := s.Init(); e != nil {
		t.Error("Can't init screen")
		return
	}

	s.Clear()

	c1Style := tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlue)

	c2Style := tcell.StyleDefault.
		Foreground(tcell.ColorBlue).
		Background(tcell.ColorGreen)

	c1 := NewView("MainComponent", s)
	c1.SetBounds(Rect{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 20,
	})
	c1.SetForegroundColor(tcell.ColorYellow)
	c1.SetBackgroundColor(tcell.ColorBlue)
	c1.SetVisible(true)

	c2 := NewView("ChildComponent", s)
	c2.SetBounds(Rect{
		X:      -2,
		Y:      -2,
		Width:  5,
		Height: 5,
	})
	c2.SetForegroundColor(tcell.ColorBlue)
	c2.SetBackgroundColor(tcell.ColorGreen)
	c2.SetParent(&c1)
	c2.SetVisible(true)

	c1.AddChild(&c2)

	c1.HandleMessage(Message{
		Handler: BroadcastHandler(),
		Type:    WmDraw,
	})

	s.Show()

	// Check component 1 ---------------------------------------------------------
	// Bottom
	for x := 1; x < 11; x++ {
		for y := 5; y < 21; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}

	// Right
	for x := 5; x < 11; x++ {
		for y := 1; y < 21; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}

	// Check component 2 ---------------------------------------------------------
	for x := 1; x < 4; x++ {
		for y := 1; y < 4; y++ {
			if e := checkCell(x, y, ' ', s, c2Style, t); e != nil {
				t.Error(e)
			}
		}
	}
}

func TestView_draw_one_view_in_another_view_full_out(t *testing.T) {
	screenStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	s := mkTestScreen(t, "")
	defer s.Fini()

	s.SetStyle(screenStyle)

	if e := s.Init(); e != nil {
		t.Error("Can't init screen")
		return
	}

	s.Clear()

	c1Style := tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlue)

	c1 := NewView("MainComponent", s)
	c1.SetBounds(Rect{
		X:      1,
		Y:      1,
		Width:  10,
		Height: 20,
	})
	c1.SetForegroundColor(tcell.ColorYellow)
	c1.SetBackgroundColor(tcell.ColorBlue)
	c1.SetVisible(true)

	c2 := NewView("ChildComponent", s)
	c2.SetBounds(Rect{
		X:      -20,
		Y:      -20,
		Width:  5,
		Height: 5,
	})
	c2.SetForegroundColor(tcell.ColorBlue)
	c2.SetBackgroundColor(tcell.ColorGreen)
	c2.SetParent(&c1)
	c2.SetVisible(true)

	c1.AddChild(&c2)

	c1.HandleMessage(Message{
		Handler: BroadcastHandler(),
		Type:    WmDraw,
	})

	s.Show()

	// Check component 1 ---------------------------------------------------------
	for x := 1; x < 11; x++ {
		for y := 1; y < 21; y++ {
			if e := checkCell(x, y, ' ', s, c1Style, t); e != nil {
				t.Error(e)
			}
		}
	}
}

func TestView_change_bound(t *testing.T) {
	screenStyle := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack)

	s := mkTestScreen(t, "")
	defer s.Fini()

	s.SetStyle(screenStyle)

	if e := s.Init(); e != nil {
		t.Error("Can't init screen")
		return
	}

	s.Clear()

	c := NewView("You know my name", s)
	c.SetBounds(Rect{
		X:      1,
		Y:      1,
		Width:  3,
		Height: 4,
	})

	newBounds := Rect{
		X:      4,
		Y:      3,
		Width:  2,
		Height: 1,
	}

	c.HandleMessage(Message{
		Handler: c.Handler(),
		Type:    WmChangeBounds,
		Value:   newBounds,
	})

	if c.GetBounds() != newBounds {
		t.Error("Bounds are not update!")
	}

	newClientBounds := Rect{
		X:      0,
		Y:      0,
		Width:  2,
		Height: 1,
	}

	if c.GetClientBounds() != newClientBounds {
		t.Error("Client bounds are not update!")
	}
}