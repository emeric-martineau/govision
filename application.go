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
	"errors"

	"github.com/gdamore/tcell"
)

// Application is base struct for create text UI.
type Application struct {
	// Main window.
	mainWindow TView
	// Quit application on Ctrl+C.
	ExitOnCtrlC bool
	// Windows list. The First item is Windows that have focus.
	windowsList *list.List
	// Message bus.
	message Bus
	// Canvas to draw.
	canvas applicationCanvas
}

// MainWindow return main windows.
func (a *Application) MainWindow() TView {
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
	if a.mainWindow == nil {
		return errors.New("main windows is nil")
	}

	a.canvas.screen.SetStyle(a.canvas.brush)

	if e := a.canvas.screen.Init(); e != nil {
		return e
	}

	return nil
}

// Run application and wait event.
func (a *Application) Run() {
	defer a.canvas.screen.Fini()

	a.canvas.screen.Clear()

	applicationHandler := ApplicationHandler()

	var msg Message
	doContinue := true

	// First time send draw message to create screen.
	a.message.Send(BuildDrawMessage(a.mainWindow.Handler()))

	go poolEvent(a.canvas.screen, a.message)

	for doContinue {
		msg = <-*a.message.Channel()

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
		}

		a.canvas.screen.Sync()
	}
}

// WindowsList return the current windows list.
// Becarefull, each call create a new array to return.
func (a *Application) WindowsList() []TView {
	wl := make([]TView, 0)

	for e := a.windowsList.Front(); e != nil; e = e.Next() {
		wl = append(wl, e.Value.(TView))
	}

	return wl
}

// Canvas return application confguration.
func (a *Application) Canvas() TCanvas {
	return &a.canvas
}

// AddWindow add window to list. If first window, she become the main window.
func (a *Application) AddWindow(w TView) {
	a.windowsList.PushFront(w)

	// TODO allow change main window
	if a.mainWindow == nil {
		a.mainWindow = w
	}
}

//------------------------------------------------------------------------------
// Internal Canvas

type applicationCanvas struct {
	brush  tcell.Style
	screen tcell.Screen
}

// SetBrush set the brush to draw.
func (a *applicationCanvas) SetBrush(b tcell.Style) {
	a.brush = b
}

// UpdateBounds call when component move or resize.
func (a *applicationCanvas) UpdateBounds(r Rect) {
}

// CreateCanvasFrom create a sub-canvas for `r` parameter.
func (a *applicationCanvas) CreateCanvasFrom(r Rect) TCanvas {
	return NewCanvas(a, r)
}

// PrintChar print a charactere.
func (a *applicationCanvas) PrintChar(x int, y int, char rune) {
	a.screen.SetContent(x, y, char, nil, a.brush)
}

// PrintCharWithBrush print a charactere with brush.
func (a *applicationCanvas) PrintCharWithBrush(x int, y int, char rune, brush tcell.Style) {
	a.screen.SetContent(x, y, char, nil, brush)
}

// Fill fill canvas zone.
func (a *applicationCanvas) Fill(bounds Rect) {
	for y := bounds.Y; y < bounds.Height; y++ {
		for x := bounds.X; x < bounds.Width; x++ {
			a.screen.SetContent(x, y, ' ', nil, a.brush)
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
			if msg.Value.(TComponent).Handler() == e.Value.(TComponent).Handler() {
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

	w.Value.(TView).HandleMessage(msg)
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
	ac := applicationCanvas{
		screen: config.Screen,
		brush: config.ScreenStyle.Style.
			Foreground(config.ScreenStyle.ForegroundColor).
			Background(config.ScreenStyle.BackgroundColor),
	}

	return Application{
		ExitOnCtrlC: true,
		windowsList: list.New(),
		message:     config.Message,
		canvas:      ac,
	}
}
