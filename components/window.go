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

// ┌─[■]    ─┐
func DrawTitleBar(titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borderStyle BorderStyle) {
	titleBar := make([]rune, base.MaxInt(titleBounds.Width, drawBounds.Width))
	titleBarLen := len(titleBar)
	indexTitleBar := 0

	titleBar[indexTitleBar] = tcell.RuneULCorner // [0]
	indexTitleBar++

	const minimumTitleBar = 7

	// Draw close only if available space for
	// ┌─[■]─┐
	if titleBarLen >= minimumTitleBar {
		titleBar[indexTitleBar] = '─' // [1]
		indexTitleBar++
		titleBar[indexTitleBar] = '[' // [2]
		indexTitleBar++
		titleBar[indexTitleBar] = '■' // [3]
		indexTitleBar++
		titleBar[indexTitleBar] = ']' // [4]
		indexTitleBar++

		// Need space before and after caption -> +2
		// ─── Title ────
		if titleBarLen > len(caption)+minimumTitleBar+2 {
			paddingLen := (titleBarLen - len(caption) - minimumTitleBar - 2)
			paddingLenLeft := paddingLen / 2

			// Border before caption
			// ──────── Hello
			for i := 0; i < paddingLenLeft; i++ {
				titleBar[indexTitleBar] = '─'
				indexTitleBar++
			}

			titleBar[indexTitleBar] = ' '
			indexTitleBar++

			c := []rune(caption)

			// Draw caption
			for i := 0; i < len(c); i++ {
				titleBar[indexTitleBar] = c[i]
				indexTitleBar++
			}

			titleBar[indexTitleBar] = ' '
			indexTitleBar++

			// Border after caption
			// Hello ────
			for ; indexTitleBar < titleBarLen-1; indexTitleBar++ {
				titleBar[indexTitleBar] = '─'
			}
		} else {

		}
	}
	//tcell.RuneULCorner
	//tcell.RuneHLine
	//[
	// X -> color ?
	//]
	// tcell.RuneHLine
	// tcell.RuneURCorner

	// TODO add up/down button

	titleBar[indexTitleBar] = tcell.RuneURCorner

	base.AppScreen().Screen().
		SetContent(titleBounds.X, titleBounds.Y, titleBar[0], titleBar[1:], style)
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
	base.Fill(clientBounds, drawBounds, style)

	// Draw title
	titleBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y,
		Width:  absoluteBounds.Width,
		Height: 1,
	}

	titleStyle := tcell.StyleDefault.
		Foreground(w.GetForegroundColor()).
		Background(tcell.ColorBlue)

	DrawTitleBar(titleBounds, drawBounds, "Hello", titleStyle, w.BorderStyle)
	/*
		titleStyle := tcell.StyleDefault.
			Foreground(w.GetForegroundColor()).
			Background(tcell.ColorBlue)

		titleBounds := base.Rect{
			X:      absoluteBounds.X,
			Y:      absoluteBounds.Y,
			Width:  absoluteBounds.Width,
			Height: 1,
		}

		base.Fill(titleBounds, drawBounds, titleStyle)
		// TODO Draw bottom
		b := []rune{'Q', 'E'}
		base.AppScreen().Screen().
			SetContent(30, 30, ' ', b, style)
	*/
	// TODO Draw left
	// TODO Draw right
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
		base.SendMessage(base.BuildDrawMessage(base.BroadcastHandler()))
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
func NewWindow(name string) Window {
	w := Window{
		View: base.NewView(name),
	}

	base.SendMessage(BuildCreateWindowMessage(&w))

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
