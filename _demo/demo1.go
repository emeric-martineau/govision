package main

// Copyright 2019 The GoVision Authors
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
	"os"
	"time"

	base "github.com/emeric-martineau/govision"
	"github.com/emeric-martineau/govision/components"

	"github.com/gdamore/tcell"
)

var timer1 components.Timer
var timer2 components.Timer
var view1 base.View
var view2 base.View
var view3 base.View

func timer2Gone(t *components.Timer) {
	v1 := view1.GetZorder()
	v2 := view2.GetZorder()

	view2.SetZorder(v1)
	view1.SetZorder(v2)

	base.SendMessage(base.BuildZorderMessage(view1.GetParent().Handler()))
}

func timer1Gone(t *components.Timer) {
	r := view3.GetBounds()

	r.Y++
	r.X++

	base.PrintStringOnScreen(
		view3.GetScreen(),
		tcell.ColorBlack,
		tcell.ColorWhite, 0, 0,
		fmt.Sprintf("%+v", r))

	base.SendMessage(base.BuildChangeBounds(view3.Handler(), r))
}

func testOnDraw(v base.TView, s tcell.Screen) {
	fmt.Println("onDraw")
}

func main() {
	rootComponent := base.NewComponent("rootComponent")
	rootComponent.SetEnabled(true)

	application, e := base.NewApplication(&rootComponent)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e = application.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s := application.GetScreen()

	view1 = base.NewView("view1", s)
	view1.SetBounds(base.Rect{
		X:      5,
		Y:      8,
		Height: 3,
		Width:  3,
	})
	view1.SetBackgroundColor(tcell.ColorBlue)
	view1.SetZorder(1)
	view1.SetEnabled(true)
	view1.SetParent(&rootComponent)
	view1.SetVisible(true)

	rootComponent.AddChild(&view1)

	view2 = base.NewView("view2", s)
	view2.SetBounds(base.Rect{
		X:      6,
		Y:      2,
		Height: 10,
		Width:  10,
	})
	view2.SetBackgroundColor(tcell.ColorRed)
	view2.SetEnabled(true)
	view2.SetParent(&rootComponent)
	view2.SetVisible(true)

	rootComponent.AddChild(&view2)

	view3 = base.NewView("view3", s)
	view3.SetBounds(base.Rect{
		X:      -5,
		Y:      -5,
		Height: 5,
		Width:  5,
	})
	view3.SetBackgroundColor(tcell.ColorYellow)
	view3.SetEnabled(true)
	view3.SetParent(&view2)
	view3.SetVisible(true)

	view2.AddChild(&view3)

	view4 := base.NewView("view4", s)
	view4.SetBounds(base.Rect{
		X:      1,
		Y:      1,
		Height: 3,
		Width:  3,
	})
	view4.SetBackgroundColor(tcell.ColorGreen)
	view4.SetEnabled(true)
	view4.SetParent(&view3)
	view4.SetVisible(true)

	view3.AddChild(&view4)

	timer1 = components.NewTimer("timer1", 1000*time.Millisecond)
	timer1.OnTimer = timer1Gone
	rootComponent.AddChild(&timer1)
	timer1.SetEnabled(true)

	timer2 = components.NewTimer("timer2", 500*time.Millisecond)
	timer2.OnTimer = timer2Gone
	rootComponent.AddChild(&timer2)
	timer2.SetEnabled(true)

	application.Run()
}
