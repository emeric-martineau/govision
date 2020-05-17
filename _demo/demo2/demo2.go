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
	appConfig := base.CreateDefaultApplicationConfig()
	application := base.NewApplication(appConfig)
	application.ShowMouseCursor = true

	mainWindow := components.NewWindow("window1", appConfig.Message, application.Canvas())

	application.AddWindow(&mainWindow)

	mainWindow.SetBounds(base.Rect{
		X:      5,
		Y:      8,
		Height: 10,
		Width:  20,
	})
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)

	window2 := components.NewWindow("window2", appConfig.Message, mainWindow.ClientCanvas())

	window2.SetBounds(base.Rect{
		X:      0,
		Y:      0,
		Height: 20,
		Width:  50,
	})
	window2.SetEnabled(true)
	window2.SetVisible(true)
	window2.SetParent(&mainWindow)
	window2.SetBackgroundColor(tcell.ColorBlue)

	mainWindow.AddChild(&window2)

	aRootWindow := components.NewWindow("window3", appConfig.Message, application.Canvas())

	application.AddWindow(&aRootWindow)

	aRootWindow.SetBounds(base.Rect{
		X:      30,
		Y:      2,
		Height: 10,
		Width:  20,
	})
	aRootWindow.SetEnabled(true)
	aRootWindow.SetVisible(true)

	if e := application.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	application.Run()
}
