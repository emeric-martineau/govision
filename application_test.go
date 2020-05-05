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

func TestApplication_Exit_on_CtrlC(t *testing.T) {
	mainWindow := NewComponent("window1")

	app, _ := NewApplication(&mainWindow)

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		s := AppScreen().Screen().(tcell.SimulationScreen)

		s.InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)

		app.Run()
	}
}

func TestApplication_Exit_on_WmQuit(t *testing.T) {
	mainWindow := NewComponent("window1")

	app, _ := NewApplication(&mainWindow)

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		// Just for code coverage
		SendMessage(BuildEmptyMessage())
		SendMessage(BuildKeyMessage(tcell.NewEventKey(tcell.KeyCtrlD, '&', tcell.ModCtrl)))
		SendMessage(BuildDrawMessage(BroadcastHandler()))

		SendMessage(Message{
			Handler: BroadcastHandler(),
			Type:    WmQuit,
		})

		app.Run()
	}
}
