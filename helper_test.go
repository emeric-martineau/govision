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
	"fmt"
	"testing"

	"github.com/gdamore/tcell"
)

func TestHelper_MaxInt(t *testing.T) {
	v := MaxInt(1, 2)

	if v != 2 {
		t.Errorf("(1) should be return 2, return %d", v)
	}

	v = MaxInt(2, 1)

	if MaxInt(2, 1) != 2 {
		t.Errorf("(2) should be return 2, return %d", v)
	}
}

func TestHelper_MinInt(t *testing.T) {
	v := MinInt(1, 2)

	if v != 1 {
		t.Errorf("(1) should be return 1, return %d", v)
	}

	v = MinInt(2, 1)

	if v != 1 {
		t.Errorf("(2) should be return 1, return %d", v)
	}
}

// Rectangle 1 in rectangle 2
func TestHelper_Intersect_Rectangle1_in_rectange2(t *testing.T) {
	r1 := Rect{
		X:      0,
		Y:      0,
		Width:  20,
		Height: 20,
	}

	r2 := Rect{
		X:      5,
		Y:      5,
		Width:  10,
		Height: 10,
	}

	result := Rect{
		X:      5,
		Y:      5,
		Width:  10,
		Height: 10,
	}

	r3 := Intersect(r1, r2)

	if r3 != result {
		t.Errorf("Should be return %+v, return %+v", result, r3)
	}
}

// Rectangle 1 intersect rectangle 2
func TestHelper_Intersect_Rectangle1_intersect_rectangle2(t *testing.T) {
	r1 := Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	}

	r2 := Rect{
		X:      5,
		Y:      5,
		Width:  15,
		Height: 15,
	}

	result := Rect{
		X:      5,
		Y:      5,
		Width:  5,
		Height: 5,
	}

	r3 := Intersect(r1, r2)

	if r3 != result {
		t.Errorf("TestIntersect2 should be return %+v, return %+v", result, r3)
	}
}

// Rectangle 1 no intersect rectangle 2
func TestHelper_Intersect_Rectangle1_no_intersect_rectangle2(t *testing.T) {
	r1 := Rect{
		X:      0,
		Y:      0,
		Width:  2,
		Height: 2,
	}

	r2 := Rect{
		X:      5,
		Y:      5,
		Width:  15,
		Height: 15,
	}

	result := Rect{
		X:      0,
		Y:      0,
		Width:  0,
		Height: 0,
	}

	r3 := Intersect(r1, r2)

	if r3 != result {
		t.Errorf("TestIntersect3 should be return %+v, return %+v", result, r3)
	}
}

// All components are visible.
func TestHelper_All_component_are_visibles_CalculateAbsolutePosition(t *testing.T) {
	s := mkTestScreen(t, "")
	view1 := NewView("view1", s)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	view2 := NewView("view2", s)
	view2.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  40,
		Height: 20,
	})
	view2.SetParent(&view1)

	view3 := NewView("view3", s)
	view3.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  30,
		Height: 10,
	})
	view3.SetParent(&view2)

	result := Rect{
		X:      6,
		Y:      9,
		Width:  30,
		Height: 10,
	}

	r := CalculateAbsolutePosition(&view3)

	if r != result {
		t.Errorf("Should be return %+v, return %+v", result, r)
	}
}

// View3 is not visible.
func TestHelper_View3_is_not_visible_CalculateAbsolutePosition2(t *testing.T) {
	s := mkTestScreen(t, "")
	view1 := NewView("view1", s)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	view2 := NewView("view2", s)
	view2.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  40,
		Height: 20,
	})
	view2.SetParent(&view1)

	view3 := NewView("view3", s)
	view3.SetBounds(Rect{
		X:      -30,
		Y:      -30,
		Width:  5,
		Height: 5,
	})
	view3.SetParent(&view2)

	result := Rect{
		X:      0,
		Y:      0,
		Width:  0,
		Height: 0,
	}

	r := CalculateAbsolutePosition(&view3)

	if r != result {
		t.Errorf("Should be return %+v, return %+v", result, r)
	}
}

// View2 is not visible.
func TestHelper_View2_is_not_visible_CalculateAbsolutePosition3(t *testing.T) {
	s := mkTestScreen(t, "")
	view1 := NewView("view1", s)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	view2 := NewView("view2", s)
	view2.SetBounds(Rect{
		X:      -100,
		Y:      -100,
		Width:  40,
		Height: 20,
	})
	view2.SetParent(&view1)

	view3 := NewView("view3", s)
	view3.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  30,
		Height: 10,
	})
	view3.SetParent(&view2)

	result := Rect{
		X:      0,
		Y:      0,
		Width:  0,
		Height: 0,
	}

	r := CalculateAbsolutePosition(&view3)

	if r != result {
		t.Errorf("Should be return %+v, return %+v", result, r)
	}
}

// With only on component.
func TestHelper_with_only_on_component_CalculateAbsolutePosition4(t *testing.T) {
	s := mkTestScreen(t, "")
	view1 := NewView("view1", s)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	result := Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	}

	r := CalculateAbsolutePosition(&view1)

	if r != result {
		t.Errorf("Should be return %+v, return %+v", result, r)
	}
}

func mkTestScreen(t *testing.T, charset string) tcell.SimulationScreen {
	s := tcell.NewSimulationScreen(charset)
	if s == nil {
		t.Fatalf("Failed to get simulation screen")
	}
	if e := s.Init(); e != nil {
		t.Fatalf("Failed to initialize screen: %v", e)
	}
	return s
}

func checkCell(x int, y int, c rune, s tcell.SimulationScreen, st tcell.Style, t *testing.T) error {
	width, _ := s.Size()

	b, _, _ := s.GetContents()

	cell := &b[y*width+x]

	if len(cell.Runes) != 1 || len(cell.Bytes) != 1 {
		return fmt.Errorf("Cell content lenght > 1 (x: %d, y: %d)", x, y)
	} else if cell.Runes[0] != c {
		return fmt.Errorf("Incorrect cell content at (x: %d, y: %d). Want '%c': Found '%c'", x, y, c, cell.Runes[0])
	} else if cell.Style != st {
		return fmt.Errorf("Incorrect style at (x: %d, y: %d). Want '%v': %v", x, y, st, cell.Style)
	}

	return nil
}

// Test if function works :)
func TestHelper_PrintStringOnScreen_function(t *testing.T) {
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

	st := tcell.StyleDefault.
		Foreground(tcell.ColorRed).
		Background(tcell.ColorBlue)

	PrintStringOnScreen(
		s,
		tcell.ColorBlue,
		tcell.ColorRed, 0, 0,
		"Hello")

	s.Show()

	if e := checkCell(0, 0, 'H', s, st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(1, 0, 'e', s, st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(2, 0, 'l', s, st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(3, 0, 'l', s, st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(4, 0, 'o', s, st, t); e != nil {
		t.Error(e)
	}
}
