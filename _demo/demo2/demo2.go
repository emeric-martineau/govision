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

	base "github.com/emeric-martineau/govision"
	"github.com/emeric-martineau/govision/components"

	"github.com/gdamore/tcell"
)

func main() {
	mainWindow := components.NewWindow("window1")
	mainWindow.SetEnabled(true)

	application, e := base.NewApplication(&mainWindow)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e = application.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	mainWindow.SetBackgroundColor(tcell.ColorYellow)
	mainWindow.SetBounds(base.Rect{
		X:      5,
		Y:      8,
		Height: 10,
		Width:  30,
	})
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)
	mainWindow.BorderStyle = components.BorderStyleSingle

	base.SendMessage(components.BuildCreateWindowMessage(&mainWindow))

	application.Run()
}
