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

func CreateTestApplicationConfig() ApplicationConfig {
	screen := tcell.NewSimulationScreen("")

	// Screen application.
	screenStyle := ApplicationStyle{
		Style:           tcell.StyleDefault,
		ForegroundColor: tcell.ColorWhite,
		BackgroundColor: tcell.ColorBlack,
	}

	return ApplicationConfig{
		ScreenStyle: screenStyle,
		Screen:      screen,
		Message:     NewBus(),
	}
}

func TestApplication_MainWindow_is_nil(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)
	e := app.Init()

	if e == nil {
		t.Error("Error should not be nil!")
	}

	if e.Error() != "main windows is nil" {
		t.Error("Error message wrong!")
	}
}

func TestApplication_Exit_on_CtrlC(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	mainWindow := NewComponent("window1", appConfig.Message)

	app := NewApplication(appConfig)

	app.AddWindow(&mainWindow)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)

		app.Run()
	}
}

func TestApplication_Exit_on_CtrlC_with_two_windows(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	mainWindow := NewComponent("window1", appConfig.Message)
	window1 := NewComponent("window2", appConfig.Message)

	app := NewApplication(appConfig)

	app.AddWindow(&mainWindow)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmCreate,
			Value:   &window1,
		})

		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmDestroy,
			Value:   &window1,
		})

		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmDestroy,
			Value:   &mainWindow,
		})

		app.Run()
	}
}

func TestApplication_Exit_on_WmQuit(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	mainWindow := NewComponent("window2", appConfig.Message)

	app := NewApplication(appConfig)

	app.AddWindow(&mainWindow)

	mainWindow.OnReceiveMessage = func(c TComponent, m Message) bool {
		if count := len(app.WindowsList()); count != 1 {
			fmt.Printf("Windows list must have only 1 component. Found %d!\n", count)
			//t.Errorf("Windows list must have only 1 component. Found %d!", count)
		}

		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmQuit,
		})

		return false
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		// Just for code coverage
		appConfig.Message.Send(BuildEmptyMessage())
		appConfig.Message.Send(BuildKeyMessage(tcell.NewEventKey(tcell.KeyCtrlD, '&', tcell.ModCtrl)))
		appConfig.Message.Send(BuildDrawMessage(BroadcastHandler()))

		app.Run()
	}
}

func TestApplication_Exit_destroy_mainwindow(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	mainWindow := NewComponent("window3", appConfig.Message)

	mainWindow.OnReceiveMessage = func(c TComponent, m Message) bool {
		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmDestroy,
			Value:   c,
		})

		return false
	}

	app := NewApplication(appConfig)

	app.AddWindow(&mainWindow)

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		appConfig.Message.Send(Message{
			Handler: mainWindow.Handler(),
			Type:    WmDraw,
		})

		app.Run()
	}

	if len(app.WindowsList()) != 0 {
		t.Errorf("Windows list are not empty! %+v\n", app.WindowsList())

		for _, w := range app.WindowsList() {
			t.Errorf("%+v\n", w)
		}
	}
}
