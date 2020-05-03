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
	"os"
	"strings"

	"github.com/gdamore/tcell"
)

// ApplicationScreen is base struct for screen of application.
type ApplicationScreen struct {
	// Default screen style.
	Style tcell.Style
	// Default screen text color.
	ForegroundColor tcell.Color
	// Default screen background color.
	BackgroundColor tcell.Color
}

// Application is base struct for create text UI.
type Application struct {
	// Main window.
	MainWindow TComponent
	// Quit application on Ctrl+C.
	ExitOnCtrlC bool
	// Screen application.
	Screen ApplicationScreen
	// Internal screen.
	screen tcell.Screen
}

// SetEncodingFallback changes the behavior of GetEncoding when a suitable
// encoding is not found.  The default is EncodingFallbackFail, which
// causes GetEncoding to simply return nil.
func (a Application) SetEncodingFallback(fb tcell.EncodingFallback) {
	tcell.SetEncodingFallback(fb)
}

// Init initialize screen (color, style...).
func (a Application) Init() error {
	style := a.Screen.Style.
		Foreground(a.Screen.ForegroundColor).
		Background(a.Screen.BackgroundColor)

	a.screen.SetStyle(style)

	if e := a.screen.Init(); e != nil {
		return e
	}

	return nil
}

// Run application and wait event.
func (a Application) Run() {
	defer a.screen.Fini()

	a.screen.Clear()

	var currentMessage Message
	doContinue := true

	// First time send draw message to create screen.
	SendMessage(BuildDrawMessage(a.MainWindow.Handler()))

	go poolEvent(a.screen)

	// Remember last message type cause many WmDraw can occure. If that, don't
	// refresh screen. Only if previous message is not a draw.
	previousMessageType := WmDraw

	for doContinue {
		currentMessage = <-busChannel

		if currentMessage.Type == WmKey && a.ExitOnCtrlC {
			if currentMessage.Value.(*tcell.EventKey).Key() == tcell.KeyCtrlC {
				doContinue = false
			} else {
				a.MainWindow.HandleMessage(currentMessage)
			}
		} else if currentMessage.Type == WmQuit {
			doContinue = false
		} else {
			a.MainWindow.HandleMessage(currentMessage)

			if currentMessage.Type == WmDraw && previousMessageType != WmDraw {
				a.screen.Sync()
			}
		}

		previousMessageType = currentMessage.Type
	}
}

// GetScreen return screen to draw.
func (a Application) GetScreen() tcell.Screen {
	return a.screen
}

// Run in go function to wait keyboard or mouse event.
func poolEvent(s tcell.Screen) {
	for {
		ev := s.PollEvent()

		switch ev := ev.(type) { // TODO Mouse message
		case *tcell.EventKey:
			SendMessage(Message{
				Handler: BroadcastHandler(),
				Type:    WmKey,
				Value:   ev,
			})
		case *tcell.EventResize:
			s.Sync()
			SendMessage(BuildScreenResizeMessage(s))
		}
	}
}

// NewApplication create a text application.
func NewApplication(mainWindow TComponent) (Application, error) {
	var s tcell.Screen
	var e error

	if strings.HasSuffix(os.Args[0], ".test") {
		s = tcell.NewSimulationScreen("")
		e = nil
	} else {
		s, e = tcell.NewScreen()
	}

	if e != nil {
		return Application{}, e
	}

	screen := ApplicationScreen{
		Style:           tcell.StyleDefault,
		ForegroundColor: tcell.ColorWhite,
		BackgroundColor: tcell.ColorBlack,
	}

	app := Application{
		MainWindow:  mainWindow,
		screen:      s,
		Screen:      screen,
		ExitOnCtrlC: true,
	}

	return app, nil
}
