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

/*
// All components are visible.
func TestHelper_All_component_are_visibles_CalculateAbsolutePosition(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	view1 := NewView("view1", appConfig)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	view2 := NewView("view2", appConfig)
	view2.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  40,
		Height: 20,
	})
	view2.SetParent(&view1)

	view3 := NewView("view3", appConfig)
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
func TestHelper_View3_is_not_visible_CalculateDrawZone2(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	view1 := NewView("view1", appConfig)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	view2 := NewView("view2", appConfig)
	view2.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  40,
		Height: 20,
	})
	view2.SetParent(&view1)

	view3 := NewView("view3", appConfig)
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

	r := CalculateDrawZone(&view3)

	if r != result {
		t.Errorf("Should be return %+v, return %+v", result, r)
	}
}

// View2 is not visible.
func TestHelper_View2_is_not_visible_CalculateDrawZone3(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	view1 := NewView("view1", appConfig)
	view1.SetBounds(Rect{
		X:      2,
		Y:      3,
		Width:  50,
		Height: 30,
	})

	view2 := NewView("view2", appConfig)
	view2.SetBounds(Rect{
		X:      -100,
		Y:      -100,
		Width:  40,
		Height: 20,
	})
	view2.SetParent(&view1)

	view3 := NewView("view3", appConfig)
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

	r := CalculateDrawZone(&view3)

	if r != result {
		t.Errorf("Should be return %+v, return %+v", result, r)
	}
}

// With only on component.
func TestHelper_with_only_on_component_CalculateAbsolutePosition4(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	view1 := NewView("view1", appConfig)
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
*/

func checkCell(screen tcell.Screen, x int, y int, c rune, st tcell.Style, t *testing.T) error {

	mainc, _, style, _ := screen.GetContent(x, y)

	//fmt.Printf("%+v\n", screen)

	if mainc != c {
		return fmt.Errorf("Incorrect cell content at (x: %d, y: %d). Want '%c': Found '%c'", x, y, c, mainc)
	} else if style != st {
		stFg, stBg, stAttr := st.Decompose()
		styleFg, styleBg, styleAttr := style.Decompose()

		return fmt.Errorf("Incorrect style at (x: %d, y: %d).\nForeground color want %v found %v\nBackground color want %v found %v\nAttributes want %v found %v", x, y, stFg, styleFg, stBg, styleBg, stAttr, styleAttr)
	}

	return nil
}

// Test if function works :)
func TestHelper_PrintStringOnScreen_function(t *testing.T) {
	appConfig := CreateTestApplicationConfig()
	defer appConfig.Screen.Fini()

	if e := appConfig.Screen.Init(); e != nil {
		t.Errorf("Can't init screen: %+v\n", e)
	}

	appConfig.Screen.Clear()

	st := tcell.StyleDefault.
		Foreground(tcell.ColorRed).
		Background(tcell.ColorBlue)

	PrintStringOnScreen(
		appConfig.Screen,
		tcell.ColorBlue,
		tcell.ColorRed, 0, 0,
		"Hello")

	appConfig.Screen.Show()

	if e := checkCell(appConfig.Screen, 0, 0, 'H', st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(appConfig.Screen, 1, 0, 'e', st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(appConfig.Screen, 2, 0, 'l', st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(appConfig.Screen, 3, 0, 'l', st, t); e != nil {
		t.Error(e)
	}
	if e := checkCell(appConfig.Screen, 4, 0, 'o', st, t); e != nil {
		t.Error(e)
	}
}
