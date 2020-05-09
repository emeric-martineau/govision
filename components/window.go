package components

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
	base "github.com/emeric-martineau/govision"
	"github.com/gdamore/tcell"
)

const (
	// BorderStyleNone windows don't have border (hard to resize.)
	BorderStyleNone = iota
	// BorderStyleSingle border with single line.
	BorderStyleSingle
	// BorderStyleDouble border with double line.
	BorderStyleDouble
	// BorderStyleEmpty border without line.
	BorderStyleEmpty
)

// BorderStyle border of window.
type BorderStyle int

// Window is the base object of all widget.
type Window struct {
	// Window title.
	Caption string
	// Border style.
	BorderStyle BorderStyle

	base.View
}

// DefaultDrawTitleBar default draw for title bar.
// Give ┌─[■]─ My title ─┐
func DefaultDrawTitleBar(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borderStyle BorderStyle) {
	indexTitleBar := 0

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneULCorner, nil, style) // [1]
	indexTitleBar++

	const minimumTitleBar = 7
	// TODO change char with border style
	// Draw close only if available space for
	// ┌─[■]─┐
	// Need space before and after caption -> +2
	if titleBounds.Width >= minimumTitleBar {
		const minimumTitleBarWithCaption = minimumTitleBar + 2

		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneHLine, nil, style) // [1]
		indexTitleBar++
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, '[', nil, style) // [2]
		indexTitleBar++
		// TODO set color for close button
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, '■', nil, style) // [3]
		indexTitleBar++
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, ']', nil, style) // [4]
		indexTitleBar++

		paddingLen := titleBounds.Width - len(caption) - 2
		paddingLenLeft := (paddingLen / 2)

		// Need space before and after caption -> +2
		// ─── Title ────
		if titleBounds.Width > minimumTitleBarWithCaption {
			// Border before caption
			// ──────── Hello
			for ; indexTitleBar < paddingLenLeft; indexTitleBar++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneHLine, nil, style)
			}

			screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, ' ', nil, style)
			indexTitleBar++

			var c []rune

			if titleBounds.Width > len(caption)+minimumTitleBarWithCaption {
				c = []rune(caption)
			} else {
				// If we don't have enought space to draw title
				len := titleBounds.Width - minimumTitleBarWithCaption
				c = []rune(caption[:len])
			}

			// Draw caption
			for i := 0; i < len(c); i++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, c[i], nil, style)
				indexTitleBar++
			}

			screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, ' ', nil, style)
			indexTitleBar++

			// Border after caption
			// Hello ────
			// -1 cause last char is corner
			for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneHLine, nil, style)
			}

			// TODO add up/down button

		} else {
			// No space to draw caption
			for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneHLine, nil, style)
			}
		}
	} else {
		// No space to draw close button
		for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
			screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneHLine, nil, style)
		}
	}

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneURCorner, nil, style)
}

// DefaultDrawBottomBar draw bottom border of window.
// Give └───────────┘
func DefaultDrawBottomBar(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borderStyle BorderStyle) {
	indexTitleBar := 0

	// TODO if not fully visible

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneLLCorner, nil, style) // [0]
	indexTitleBar++

	// TODO change char with border style
	// -1 cause last char is corner
	for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneHLine, nil, style)
	}

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, tcell.RuneLRCorner, nil, style)
}

// DefaultDrawLeftOrRightBorder draw left or right border of window.
func DefaultDrawLeftOrRightBorder(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borderStyle BorderStyle) {
	titleBarLen := titleBounds.Height

	for indexTitleBar := 0; indexTitleBar < titleBarLen; indexTitleBar++ {
		screen.SetContent(titleBounds.X, titleBounds.Y+indexTitleBar, tcell.RuneVLine, nil, style)
	}
}

// GetClientBounds return client bounds.
func (w Window) GetClientBounds() base.Rect {
	return calculateClientBounds(w.GetBounds(), w.BorderStyle)
}

// Draw the view.
func (w Window) Draw() {
	if !w.GetVisible() {
		return
	}

	style := tcell.StyleDefault.
		Foreground(w.GetForegroundColor()).
		Background(w.GetBackgroundColor())

	// Get parent X and Y
	absoluteBounds := base.CalculateAbsolutePosition(&w)
	// Get real zone whe can draw
	drawBounds := base.CalculateDrawZone(&w)
	// Get client bould to draw
	clientBounds := calculateClientBounds(absoluteBounds, w.BorderStyle)

	// Draw background of window
	base.Fill(w.AppConfig().Screen, clientBounds, drawBounds, style)

	borderStyle := tcell.StyleDefault.
		Foreground(w.GetForegroundColor()).
		Background(tcell.ColorBlue)

	drawTitle(absoluteBounds, drawBounds, &w, borderStyle)

	drawBottom(absoluteBounds, drawBounds, &w, borderStyle)

	drawBorderLeft(absoluteBounds, drawBounds, &w, borderStyle)

	drawBorderRight(absoluteBounds, drawBounds, &w, borderStyle)
}

func drawTitle(absoluteBounds base.Rect, drawBounds base.Rect, w *Window, borderStyle tcell.Style) {
	titleBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y,
		Width:  absoluteBounds.Width,
		Height: 1,
	}

	DefaultDrawTitleBar(w.AppConfig().Screen, titleBounds, drawBounds, w.Caption,
		borderStyle, w.BorderStyle)
}

func drawBottom(absoluteBounds base.Rect, drawBounds base.Rect, w *Window, borderStyle tcell.Style) {
	bottomBorderBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y + absoluteBounds.Height - 1,
		Width:  absoluteBounds.Width,
		Height: 1,
	}

	DefaultDrawBottomBar(w.AppConfig().Screen, bottomBorderBounds, drawBounds, w.Caption,
		borderStyle, w.BorderStyle)
}

func drawBorderLeft(absoluteBounds base.Rect, drawBounds base.Rect, w *Window, borderStyle tcell.Style) {
	leftBorderBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y + 1,
		Width:  1,
		Height: absoluteBounds.Height - 2,
	}

	DefaultDrawLeftOrRightBorder(w.AppConfig().Screen, leftBorderBounds, drawBounds, w.Caption,
		borderStyle, w.BorderStyle)
}

func drawBorderRight(absoluteBounds base.Rect, drawBounds base.Rect, w *Window, borderStyle tcell.Style) {
	leftBorderBounds := base.Rect{
		X:      absoluteBounds.X + absoluteBounds.Width - 1,
		Y:      absoluteBounds.Y + 1,
		Width:  1,
		Height: absoluteBounds.Height - 2,
	}

	DefaultDrawLeftOrRightBorder(w.AppConfig().Screen, leftBorderBounds, drawBounds, w.Caption,
		borderStyle, w.BorderStyle)
}

// Manage message if it's for me.
// Return true to stop message propagation.
func (w Window) manageMyMessage(msg base.Message) {
	switch msg.Type {
	case base.WmDraw:
		if w.OnDraw != nil {
			w.OnDraw(&w)
		} else {
			w.Draw()
			// Redraw children.
			w.Component.HandleMessage(base.BuildDrawMessage(base.BroadcastHandler()))
		}
	case base.WmChangeBounds:
		// TODO minimum Width/Height -> 2
		w.SetBounds(msg.Value.(base.Rect))
		// Redraw all components cause maybe overide a component with Zorder
		w.AppConfig().Message.Send(base.BuildDrawMessage(base.BroadcastHandler()))
	}
}

// HandleMessage is use to manage message.
func (w Window) HandleMessage(msg base.Message) bool {

	switch msg.Handler {
	case w.Handler():
		w.manageMyMessage(msg)
		return true
	case base.BroadcastHandler():
		w.manageMyMessage(msg)
	}

	// Because Component send message to child if broadcast or draw.
	return w.Component.HandleMessage(base.BuildDrawMessage(base.BroadcastHandler()))
}

func calculateClientBounds(bounds base.Rect, borderStyle BorderStyle) base.Rect {
	switch borderStyle {
	case BorderStyleNone:
		// Remove titlebar
		bounds.Y++
		bounds.Height--
	default:
		// Remove titlebar and border
		bounds.Y++
		bounds.X++
		bounds.Height -= 2
		bounds.Width -= 2
	}

	return bounds
}

// NewWindow create new window.
func NewWindow(name string, config base.ApplicationConfig) Window {
	w := Window{
		View:    base.NewView(name, config),
		Caption: name,
	}

	return w
}

// BuildCreateWindowMessage return a message when create window.
func BuildCreateWindowMessage(w *Window) base.Message {
	return base.Message{
		Handler: base.ApplicationHandler(),
		Type:    base.WmCreate,
		Value:   w,
	}
}

// BuildDestroyWindowMessage return a message when destroy window.
func BuildDestroyWindowMessage(w *Window) base.Message {
	return base.Message{
		Handler: base.ApplicationHandler(),
		Type:    base.WmDestroy,
		Value:   w,
	}
}
