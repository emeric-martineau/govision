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
	"errors"
	"fmt"
	"testing"

	"github.com/gdamore/tcell"
)

type ErrorScreen struct {
	tcell.Screen
}

func (e ErrorScreen) SetStyle(style tcell.Style) {

}

func (e ErrorScreen) Init() error {
	return errors.New("Fake error")
}

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

func TestApplication_Error_on_screen_init(t *testing.T) {
	appConfig := CreateTestApplicationConfig()
	appConfig.Screen = ErrorScreen{}

	app := NewApplication(appConfig)

	v := NewView("name", appConfig.Message, app.Canvas())

	app.AddWindow(&v)

	e := app.Init()

	if e == nil {
		t.Error("Error should not be nil!")
	}

	if e != nil && e.Error() != "Fake error" {
		t.Error("Error should be nil!")
	}
}

func TestApplication_MainWindow_is_nil(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)
	e := app.Init()

	if e != nil && e.Error() != "main windows is nil" {
		t.Error("Error should be nil!")
	}
}

func TestApplication_Exit_on_CtrlC(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)

	mainWindow := NewView("window1", appConfig.Message, app.Canvas())

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

	app := NewApplication(appConfig)

	mainWindow := NewView("window2", appConfig.Message, app.Canvas())
	window1 := NewView("window3", appConfig.Message, app.Canvas())

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

	app := NewApplication(appConfig)

	mainWindow := NewView("window4", appConfig.Message, app.Canvas())

	app.AddWindow(&mainWindow)

	mainWindow.SetOnReceiveMessage(func(c TComponent, m Message) bool {
		if count := len(app.WindowsList()); count != 1 {
			fmt.Printf("Windows list must have only 1 component. Found %d!\n", count)
			//t.Errorf("Windows list must have only 1 component. Found %d!", count)
		}

		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmQuit,
		})

		return false
	})

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

	app := NewApplication(appConfig)

	mainWindow := NewView("window5", appConfig.Message, app.Canvas())

	mainWindow.SetOnReceiveMessage(func(c TComponent, m Message) bool {
		appConfig.Message.Send(Message{
			Handler: ApplicationHandler(),
			Type:    WmDestroy,
			Value:   &mainWindow,
		})

		return false
	})

	app.AddWindow(&mainWindow)

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Run()
	}

	if len(app.WindowsList()) != 0 {
		t.Errorf("Windows list are not empty! %+v\n", app.WindowsList())

		for _, w := range app.WindowsList() {
			t.Errorf("%+v\n", w)
		}
	}
}

func TestApplication_Canvas_draw(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)

	mainWindow := NewView("window6", appConfig.Message, app.Canvas())

	app.AddWindow(&mainWindow)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)

		canvas := app.Canvas()

		st1 := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
		canvas.SetBrush(st1)
		canvas.PrintChar(0, 0, '&')

		st2 := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
		canvas.PrintCharWithBrush(1, 1, '!', st2)

		st3 := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
		r3 := Rect{
			X:      2,
			Y:      2,
			Height: 3,
			Width:  4,
		}
		canvas.SetBrush(st3)
		canvas.Fill(r3)

		appConfig.Screen.Show()

		checkCell(appConfig.Screen, 0, 0, '&', st1, t)
		checkCell(appConfig.Screen, 0, 0, '!', st2, t)

		for x := r3.X; x < r3.X+r3.Width; x++ {
			for y := r3.Y; y < r3.Y+r3.Height; y++ {
				if e := checkCell(appConfig.Screen, x, y, ' ', st3, t); e != nil {
					t.Error(e)
				}
			}
		}

		// Just for code coverage
		canvas.UpdateBounds(r3)

		appConfig.Screen.Fini()
	}
}

// Why not working ?
func TestApplication_Event_resize(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)

	mainWindow := NewView("window7", appConfig.Message, app.Canvas())

	app.AddWindow(&mainWindow)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).SetSize(10, 20)
		appConfig.Screen.Show()
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)

		app.Run()
	}
}

func eventMouseClickChangeFocus(t *testing.T, x, y int, button tcell.ButtonMask) (TView, TView) {
	appConfig := CreateTestApplicationConfig()
	mainWindowIsCalled := false
	window2IsCalled := false

	app := NewApplication(appConfig)

	mainWindow := NewView("window8", appConfig.Message, app.Canvas())
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)
	mainWindow.SetBounds(Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	})
	mainWindow.SetOnActivate(func(s bool) {
		mainWindowIsCalled = true

		if window2IsCalled {
			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}
	})

	app.AddWindow(&mainWindow)

	window2 := NewView("window9", appConfig.Message, app.Canvas())
	window2.SetEnabled(true)
	window2.SetVisible(true)
	window2.SetBounds(Rect{
		X:      10,
		Y:      10,
		Width:  10,
		Height: 10,
	})
	window2.SetOnActivate(func(s bool) {
		window2IsCalled = true

		if mainWindowIsCalled {
			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}
	})

	app.AddWindow(&window2)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(x, y, button, 0)
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(x, y, tcell.ButtonNone, 0)

		app.Run()
	}

	return &mainWindow, &window2
}

func TestApplication_Event_mouse_click_right_change_focus(t *testing.T) {
	mainWindow, window2 := eventMouseClickChangeFocus(t, 5, 5, tcell.Button3)

	if window2.GetFocused() {
		t.Error("Window2 must be inactive")
	}

	if !mainWindow.GetFocused() {
		t.Error("Main window must be active")
	}
}

func TestApplication_Event_mouse_click_left_change_focus(t *testing.T) {
	mainWindow, window2 := eventMouseClickChangeFocus(t, 5, 5, tcell.Button1)

	if window2.GetFocused() {
		t.Error("Window2 must be inactive")
	}

	if !mainWindow.GetFocused() {
		t.Error("Main window must be active")
	}
}

func TestApplication_Event_mouse_click_out_of_window(t *testing.T) {
	appConfig := CreateTestApplicationConfig()
	mainWindowIsCalled := false
	window2IsCalled := false

	app := NewApplication(appConfig)

	mainWindow := NewView("window10", appConfig.Message, app.Canvas())
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)
	mainWindow.SetBounds(Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	})
	mainWindow.SetOnActivate(func(s bool) {
		mainWindowIsCalled = true

		if window2IsCalled {
			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}
	})

	app.AddWindow(&mainWindow)

	window2 := NewView("window11", appConfig.Message, app.Canvas())
	window2.SetEnabled(true)
	window2.SetVisible(true)
	window2.SetBounds(Rect{
		X:      10,
		Y:      10,
		Width:  10,
		Height: 10,
	})
	window2.SetOnActivate(func(s bool) {
		window2IsCalled = true

		if mainWindowIsCalled {
			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}
	})

	app.AddWindow(&window2)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(50, 50, tcell.Button1, 0)
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(50, 50, tcell.ButtonNone, 0)
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)

		app.Run()
	}
}

func TestApplication_Event_mouse_click_out(t *testing.T) {
	isMouseClick := false
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)
	app.ShowMouseCursor = true

	mainWindow := NewView("window12", appConfig.Message, app.Canvas())
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)
	mainWindow.SetBounds(Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	})
	mainWindow.SetOnReceiveMessage(func(c TComponent, msg Message) bool {
		fmt.Printf("%+v\n", msg)
		if msg.Type == WmLButtonDown {
			isMouseClick = true

			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}

		return false
	})
	app.AddWindow(&mainWindow)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(5, 5, tcell.Button1, 0)

		app.Run()
	}

	if !isMouseClick {
		t.Error("No WmLButtonDown event receive")
	}
}

func TestApplication_Mouse_enter_mouse_leave(t *testing.T) {
	isMouseEnter := false
	isMouseLeave := false
	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)

	mainWindow := NewView("window13", appConfig.Message, app.Canvas())
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)
	mainWindow.SetBounds(Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	})
	mainWindow.SetOnReceiveMessage(func(c TComponent, msg Message) bool {
		if msg.Type == WmMouseEnter {
			isMouseEnter = true
		} else if msg.Type == WmMouseLeave {
			isMouseLeave = true

			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}

		return false
	})

	app.AddWindow(&mainWindow)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(5, 5, tcell.ButtonNone, 0)
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(10, 10, tcell.ButtonNone, 0)

		app.Run()

		if !isMouseEnter {
			t.Error("Error no mouse enter event")
		}

		if !isMouseLeave {
			t.Error("Error no mouse enter event")
		}
	}
}

func TestApplication_Mouse_enter_mouse_leave_with_multi_window(t *testing.T) {
	isMouseEnterMainWindow := false
	isMouseLeaveMainWindow := false
	isMouseEnterWindow2 := false

	appConfig := CreateTestApplicationConfig()

	app := NewApplication(appConfig)

	mainWindow := NewView("window14", appConfig.Message, app.Canvas())
	mainWindow.SetEnabled(true)
	mainWindow.SetVisible(true)
	mainWindow.SetBounds(Rect{
		X:      0,
		Y:      0,
		Width:  10,
		Height: 10,
	})
	mainWindow.SetOnReceiveMessage(func(c TComponent, msg Message) bool {
		if msg.Type == WmMouseEnter {
			isMouseEnterMainWindow = true
		} else if msg.Type == WmMouseLeave {
			isMouseLeaveMainWindow = true
		}

		return false
	})

	app.AddWindow(&mainWindow)

	window2 := NewView("window15", appConfig.Message, app.Canvas())
	window2.SetEnabled(true)
	window2.SetVisible(true)
	window2.SetBounds(Rect{
		X:      10,
		Y:      10,
		Width:  10,
		Height: 10,
	})
	window2.SetOnReceiveMessage(func(c TComponent, msg Message) bool {
		if msg.Type == WmMouseEnter {
			isMouseEnterWindow2 = true
			app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectKey(tcell.KeyCtrlC, ' ', tcell.ModCtrl)
		}

		return false
	})

	app.AddWindow(&window2)

	if app.MainWindow() != &mainWindow {
		t.Error("Error main window are different")
	}

	app.SetEncodingFallback(tcell.EncodingFallbackASCII)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(5, 5, tcell.ButtonNone, 0)
		app.Canvas().(*applicationCanvas).screen.(tcell.SimulationScreen).InjectMouse(11, 11, tcell.ButtonNone, 0)

		app.Run()

		if !isMouseEnterMainWindow {
			t.Error("Error no mouse enter event")
		}

		if !isMouseLeaveMainWindow {
			t.Error("Error no mouse enter event")
		}

		if !isMouseEnterWindow2 {
			t.Error("Error no mouse enter event (windows2)")
		}
	}
}
