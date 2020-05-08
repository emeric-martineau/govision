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
	"os"
	"strings"

	"container/list"

	"github.com/gdamore/tcell"
)

// Application is base struct for create text UI.
type Application struct {
	// Main window.
	mainWindow TComponent
	// Quit application on Ctrl+C.
	ExitOnCtrlC bool
	// Windows list. The First item is Windows that have focus.
	windowsList *list.List
}

// MainWindow return main windows.
func (a Application) MainWindow() TComponent {
	return a.mainWindow
}

// SetEncodingFallback changes the behavior of GetEncoding when a suitable
// encoding is not found.  The default is EncodingFallbackFail, which
// causes GetEncoding to simply return nil.
func (a Application) SetEncodingFallback(fb tcell.EncodingFallback) {
	tcell.SetEncodingFallback(fb)
}

// Init initialize screen (color, style...).
func (a Application) Init() error {
	style := appScreen.Style.
		Foreground(appScreen.ForegroundColor).
		Background(appScreen.BackgroundColor)

	screen.SetStyle(style)

	if e := screen.Init(); e != nil {
		return e
	}

	return nil
}

// Run application and wait event.
func (a Application) Run() {
	defer screen.Fini()

	screen.Clear()

	applicationHandler := ApplicationHandler()

	var msg Message
	doContinue := true

	// First time send draw message to create screen.
	SendMessage(BuildDrawMessage(a.mainWindow.Handler()))

	go poolEvent(screen)

	// Remember last message type cause many WmDraw can occure. If that, don't
	// refresh screen. Only if previous message is not a draw.
	previousMessageType := WmDraw

	for doContinue {
		msg = <-busChannel

		if msg.Type == WmKey && a.ExitOnCtrlC {
			if msg.Value.(*tcell.EventKey).Key() == tcell.KeyCtrlC {
				doContinue = false
			} else {
				doContinue = a.callFocusedWindowHandleMessage(msg)
			}
		} else if msg.Handler == applicationHandler {
			doContinue = a.manageMyMessage(msg)
		} else {
			doContinue = a.callFocusedWindowHandleMessage(msg)

			if msg.Type == WmDraw && previousMessageType != WmDraw {
				screen.Sync()
			}
		}

		previousMessageType = msg.Type
	}
}

// WindowsList return the current windows list.
// Becarefull, each call create a new array to return.
func (a Application) WindowsList() []TComponent {
	wl := make([]TComponent, 0)

	for e := a.windowsList.Front(); e != nil; e = e.Next() {
		wl = append(wl, e.Value.(TComponent))
	}

	return wl
}

func (a Application) manageMyMessage(msg Message) bool {
	switch msg.Type {
	case WmQuit:
		return false
	case WmCreate:
		// Add window to list
		a.windowsList.PushFront(msg.Value)
	case WmDestroy:
		// Remove window to list and check is MainWindow
		for e := a.windowsList.Front(); e != nil; e = e.Next() {
			fmt.Printf("%+v\n", e.Value)
			if msg.Value == e.Value {
				a.windowsList.Remove(e)

				if msg.Value == a.mainWindow {
					return false
				}

				return true
			}
		}
	}

	return true
}

// Call focused windows if not nil.
func (a Application) callFocusedWindowHandleMessage(msg Message) bool {
	w := a.windowsList.Front()

	if w == nil {
		return false
	}

	w.Value.(TComponent).HandleMessage(msg)

	return true
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
			SendMessage(BuildScreenResizeMessage())
		}
	}
}

// NewApplication create a text application.
func NewApplication(mainWindow TComponent) (Application, error) {
	var e error

	if strings.HasSuffix(os.Args[0], ".test") {
		screen = tcell.NewSimulationScreen("")
		e = nil
	} else {
		screen, e = tcell.NewScreen()
	}

	if e != nil {
		return Application{}, e
	}

	app := Application{
		mainWindow:  mainWindow,
		ExitOnCtrlC: true,
		windowsList: list.New(),
	}

	return app, nil
}
