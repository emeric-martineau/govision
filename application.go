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
	"github.com/google/uuid"
)

type lastCursorPosAndStyle struct {
	x     int
	y     int
	mainc rune
	combc []rune
	style tcell.Style
}

// Application is base struct for create text UI.
type Application struct {
	// Main window.
	mainWindow TView
	// Remember last windows under cursor to sent mouse move, enter, leave.
	lastWindowUnderMouse TView
	// Quit application on Ctrl+C.
	ExitOnCtrlC bool
	// Show text mouse cursor.
	ShowMouseCursor bool
	// Windows list. The First item is Windows that have focus.
	windowsList *list.List
	// Message bus.
	message Bus
	// Canvas to draw.
	canvas applicationCanvas
	// Keep previous event of mouse to know if cursor move, click...
	previousMousEvent tcell.EventMouse
	// Last cursor.
	lastCursorPosAndStyle lastCursorPosAndStyle
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

	a.canvas.screen.EnableMouse()

	return nil
}

// Run application and wait event.
func (a *Application) Run() {
	defer a.canvas.screen.Fini()

	a.canvas.screen.Clear()

	a.storeCursorInfo(0, 0)

	applicationHandler := ApplicationHandler()

	var msg Message
	doContinue := true

	// First time send draw message to create screen.
	a.message.Send(BuildDrawMessage(applicationHandler))

	go poolEvent(a.canvas.screen, a.message)

	a.windowsList.Front().Value.(TView).SetFocused(true)

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
			a.callWindowHandleMessage(msg)
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
	for y := bounds.Y; y < bounds.Y+bounds.Height; y++ {
		for x := bounds.X; x < bounds.X+bounds.Width; x++ {
			a.screen.SetContent(x, y, ' ', nil, a.brush)
		}
	}
}

//------------------------------------------------------------------------------
// Internal functions

func (a *Application) storeCursorInfo(x, y int) {
	a.lastCursorPosAndStyle.x = x
	a.lastCursorPosAndStyle.y = y

	mainc, combc, style, _ := a.canvas.screen.GetContent(x, y)

	a.lastCursorPosAndStyle.mainc = mainc
	a.lastCursorPosAndStyle.combc = combc
	a.lastCursorPosAndStyle.style = style
}

func (a *Application) manageMyMessage(msg Message) bool {
	switch msg.Type {
	case WmMouse:
		a.manageMouseMessage(msg)
	case WmDraw:
		for e := a.windowsList.Front(); e != nil; e = e.Next() {
			a.message.Send(BuildDrawMessage(BroadcastHandler()))
		}
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

func (a *Application) findWindowsByCoordinate(x int, y int) (*list.Element, TView) {
	var currentWindow TView
	var currentWindowBounds Rect

	for e := a.windowsList.Front(); e != nil; e = e.Next() {
		currentWindow = e.Value.(TView)
		currentWindowBounds = currentWindow.GetBounds()

		if currentWindow.GetVisible() &&
			InVertical(y, currentWindowBounds) &&
			InHorizontal(x, currentWindowBounds) {
			return e, currentWindow
		}
	}

	return nil, nil
}

func (a *Application) findWindowsByHandle(handle uuid.UUID) (*list.Element, TView) {
	var currentWindow TView

	for e := a.windowsList.Front(); e != nil; e = e.Next() {
		currentWindow = e.Value.(TView)
		if currentWindow.Handler() == handle {
			return e, currentWindow
		}
	}

	return nil, nil
}

func (a *Application) manageMouseClickDown(ev *tcell.EventMouse, side uint) {
	// Ok send event
	x, y := ev.Position()

	e, window := a.findWindowsByCoordinate(x, y)

	if window == nil {
		return
	}

	// Is windows has already focus ?
	currentFocusedWindow := a.windowsList.Front().Value.(TView)

	if window.Handler() == currentFocusedWindow.Handler() {
		// Send a click message
		a.message.Send(BuildClickMouseMessage(window.Handler(), ev, side))
	} else {
		// Send focus message
		a.windowsList.MoveToFront(e)

		a.message.Send(BuildDesactivateMessage(currentFocusedWindow.Handler()))
		a.message.Send(BuildActivateMessage(window.Handler()))
	}
}

func (a *Application) manageMouseClickUp(ev *tcell.EventMouse, side uint) {
	// Ok send event
	x, y := ev.Position()

	_, window := a.findWindowsByCoordinate(x, y)

	if window == nil {
		return
	}

	a.message.Send(BuildClickMouseMessage(window.Handler(), ev, side))
}

func (a *Application) displayMouseCursor(x, y int) {
	// TODO Get cursor type of window to draw cross, hour glass...
	if a.ShowMouseCursor && (a.lastCursorPosAndStyle.x != x || a.lastCursorPosAndStyle.y != y) {
		// Cursor move
		// First restore last position
		a.canvas.screen.SetContent(
			a.lastCursorPosAndStyle.x,
			a.lastCursorPosAndStyle.y,
			a.lastCursorPosAndStyle.mainc,
			a.lastCursorPosAndStyle.combc,
			a.lastCursorPosAndStyle.style)

		// Save current data
		a.storeCursorInfo(x, y)

		st := a.lastCursorPosAndStyle.style.Reverse(true)

		a.canvas.screen.SetContent(
			a.lastCursorPosAndStyle.x,
			a.lastCursorPosAndStyle.y,
			a.lastCursorPosAndStyle.mainc,
			a.lastCursorPosAndStyle.combc,
			st)
	}
}

func (a *Application) manageMouseMessage(msg Message) {
	ev := msg.Value.(*tcell.EventMouse)
	checkMouseMove := true

	// Left click
	// Check if before left click is active
	if ev.Buttons()&tcell.Button1 != 0 && a.previousMousEvent.Buttons()&tcell.Button1 == 0 {
		a.manageMouseClickDown(ev, WmLButtonDown)
		checkMouseMove = false
	} else if ev.Buttons()&tcell.Button1 == 0 && a.previousMousEvent.Buttons()&tcell.Button1 != 0 {
		a.manageMouseClickUp(ev, WmLButtonUp)
		checkMouseMove = false
	}

	// Right click
	// Check if before right click is active
	if ev.Buttons()&tcell.Button3 != 0 && a.previousMousEvent.Buttons()&tcell.Button3 == 0 {
		a.manageMouseClickDown(ev, WmRButtonDown)
		checkMouseMove = false
	} else if ev.Buttons()&tcell.Button3 == 0 && a.previousMousEvent.Buttons()&tcell.Button3 != 0 {
		a.manageMouseClickUp(ev, WmRButtonUp)
		checkMouseMove = false
	}

	x, y := ev.Position()

	// Check mouse move only if not click message send
	if checkMouseMove {
		_, window := a.findWindowsByCoordinate(x, y)

		if window == nil {
			// Send mouse leave
			if a.lastWindowUnderMouse != nil {
				a.message.Send(BuildMouseLeaveMessage(a.lastWindowUnderMouse.Handler()))
				a.lastWindowUnderMouse = nil
			}
		} else {
			if a.lastWindowUnderMouse == nil {
				// Send mouve enter
				a.message.Send(BuildMouseEnterMessage(window.Handler(), x, y))
			} else if a.lastWindowUnderMouse != nil && a.lastWindowUnderMouse.Handler() == window.Handler() {
				// Mouse move
			} else {
				// Send mouse leave
				a.message.Send(BuildMouseLeaveMessage(a.lastWindowUnderMouse.Handler()))

				// Send mouve enter
				a.message.Send(BuildMouseEnterMessage(window.Handler(), x, y))
			}

			a.lastWindowUnderMouse = window
		}
	}

	a.displayMouseCursor(x, y)

	a.previousMousEvent = *ev
}

// Call focused windows if not nil.
func (a *Application) callFocusedWindowHandleMessage(msg Message) {
	w := a.windowsList.Front()

	w.Value.(TView).HandleMessage(msg)
}

// Call windows by handle.
func (a *Application) callWindowHandleMessage(msg Message) {
	if msg.Handler == BroadcastHandler() {
		var currentWindow TView

		for e := a.windowsList.Front(); e != nil; e = e.Next() {
			currentWindow = e.Value.(TView)
			currentWindow.HandleMessage(msg)
		}
	} else {
		_, w := a.findWindowsByHandle(msg.Handler)

		w.HandleMessage(msg)
	}
}

// Run in go function to wait keyboard or mouse event.
func poolEvent(screen tcell.Screen, message Bus) {
	for {
		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventMouse:
			message.Send(Message{
				Handler: ApplicationHandler(),
				Type:    WmMouse,
				Value:   ev,
			})

			screen.Sync()
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
