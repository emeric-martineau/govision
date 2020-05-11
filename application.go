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
	// Application config
	config ApplicationConfig
}

// MainWindow return main windows.
func (a *Application) MainWindow() TComponent {
	return a.mainWindow
}

// SetEncodingFallback changes the behavior of GetEncoding when a suitable
// encoding is not found.  The default is EncodingFallbackFail, which
// causes GetEncoding to simply return nil.
func (a *Application) SetEncodingFallback(fb tcell.EncodingFallback) {
	tcell.SetEncodingFallback(fb)
}

// Init initialize screen (color, style...).
func (a *Application) Init() error {
	style := a.config.ScreenStyle.Style.
		Foreground(a.config.ScreenStyle.ForegroundColor).
		Background(a.config.ScreenStyle.BackgroundColor)

	a.config.Screen.SetStyle(style)

	if e := a.config.Screen.Init(); e != nil {
		return e
	}

	return nil
}

// Run application and wait event.
func (a *Application) Run() {
	defer a.config.Screen.Fini()

	a.config.Screen.Clear()

	applicationHandler := ApplicationHandler()

	var msg Message
	doContinue := true

	// First time send draw message to create screen.
	a.config.Message.Send(BuildDrawMessage(a.mainWindow.Handler()))

	go poolEvent(a.config.Screen, a.config.Message)

	// Remember last message type cause many WmDraw can occure. If that, don't
	// refresh screen. Only if previous message is not a draw.
	previousMessageType := WmDraw

	for doContinue {
		msg = <-*a.config.Message.Channel()

		if msg.Type == WmKey && a.ExitOnCtrlC {
			if msg.Value.(*tcell.EventKey).Key() == tcell.KeyCtrlC {
				doContinue = false
			} else {
				a.callFocusedWindowHandleMessage(msg)
			}
		} else if msg.Handler == applicationHandler {
			doContinue = a.manageMyMessage(msg)
		} else {
			a.callFocusedWindowHandleMessage(msg)

			if msg.Type == WmDraw && previousMessageType != WmDraw {
				a.config.Screen.Sync()
			}
		}

		previousMessageType = msg.Type
	}
}

// WindowsList return the current windows list.
// Becarefull, each call create a new array to return.
func (a *Application) WindowsList() []TComponent {
	wl := make([]TComponent, 0)

	for e := a.windowsList.Front(); e != nil; e = e.Next() {
		wl = append(wl, e.Value.(TComponent))
	}

	return wl
}

// Config return application confguration.
func (a *Application) Config() ApplicationConfig {
	return a.config
}

// AddWindow add window to list. If first window, she become the main window.
func (a *Application) AddWindow(w TComponent) {
	a.windowsList.PushFront(w)

	// TODO allow change main window
	if a.mainWindow == nil {
		a.mainWindow = w
	}
}

//------------------------------------------------------------------------------
// Internal Canvas

// TODO is too strange to use application for canvas. Add method Canvas()

// SetBrush set the brush to draw.
func (a *Application) SetBrush(b tcell.Style) {
	a.config.ScreenStyle.Style = b
}

// UpdateBounds call when component move or resize.
func (a *Application) UpdateBounds(r Rect) {
}

// CreateCanvasFrom create a sub-canvas for `r` parameter.
func (a *Application) CreateCanvasFrom(r Rect) TCanvas {
	return NewCanvas(a, r)
}

// PrintChar print a charactere.
func (a *Application) PrintChar(x int, y int, char rune) {
	a.config.Screen.SetContent(x, y, char, nil, a.config.ScreenStyle.Style)
}

// PrintCharWithBrush print a charactere with brush.
func (a *Application) PrintCharWithBrush(x int, y int, char rune, brush tcell.Style) {
	a.config.Screen.SetContent(x, y, char, nil, brush)
}

// Fill fill canvas zone.
func (a *Application) Fill(bounds Rect) {
	for y := bounds.Y; y < bounds.Height; y++ {
		for x := bounds.X; x < bounds.Width; x++ {
			a.config.Screen.SetContent(x, y, ' ', nil, a.config.ScreenStyle.Style)
		}
	}
}

//------------------------------------------------------------------------------
// Internal function

func (a *Application) manageMyMessage(msg Message) bool {
	switch msg.Type {
	case WmQuit:
		return false
	case WmCreate:
		// Add window to list
		a.windowsList.PushFront(msg.Value)
	case WmDestroy:
		// Remove window to list and check is MainWindow
		for e := a.windowsList.Front(); e != nil; e = e.Next() {
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
// Return false if need stop application.
func (a *Application) callFocusedWindowHandleMessage(msg Message) {
	w := a.windowsList.Front()

	w.Value.(TComponent).HandleMessage(msg)
}

// Run in go function to wait keyboard or mouse event.
func poolEvent(screen tcell.Screen, message Bus) {
	for {
		ev := screen.PollEvent()

		switch ev := ev.(type) { // TODO Mouse message
		case *tcell.EventKey:
			message.Send(Message{
				Handler: BroadcastHandler(),
				Type:    WmKey,
				Value:   ev,
			})
		case *tcell.EventResize:
			screen.Sync()
			message.Send(BuildScreenResizeMessage(screen))
		}
	}
}

//------------------------------------------------------------------------------
// Constructor.

// NewApplication create a text application.
func NewApplication(config ApplicationConfig) Application {
	return Application{
		ExitOnCtrlC: true,
		windowsList: list.New(),
		config:      config,
	}
}
