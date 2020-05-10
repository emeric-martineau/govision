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
	// BorderTypeSingle border with single line.
	BorderTypeSingle = 0
	// BorderTypeDouble border with double line.
	BorderTypeDouble = 1
	// BorderTypeEmpty border without line.
	BorderTypeEmpty = 2

	// Index for draw windows border.

	// ULCorner upper left corner.
	ULCorner = 0
	// HLine horizontal line.
	HLine = 1
	// CloseLeft close character left.
	CloseLeft = 2
	// CloseRight close character right.
	CloseRight = 3
	// Close close character.
	Close = 4
	// URCorner upper right corner.
	URCorner = 5
	// CaptionSpace space before/after caption in title bar.
	CaptionSpace = 6
	// LLCorner lower left corner.
	LLCorner = 7
	// LRCorner lower right corner.
	LRCorner = 8
	// VLine vertical line.
	VLine = 9
)

var bordersChars [][]rune

// BorderType border of window.
type BorderType int

// WindowBorder is border style
type WindowBorder struct {
	// Border type.
	Type BorderType
	// Border color.
	BackgroundColor tcell.Color
	// Border color.
	ForegroundColor tcell.Color
}

// Window is the base object of all widget.
type Window struct {
	// Window title.
	Caption string
	// Border
	Border WindowBorder

	base.View
}

func init() {
	bordersChars = make([][]rune, 3)

	borderTypeSingle := make([]rune, 10)

	borderTypeSingle[ULCorner] = tcell.RuneULCorner
	borderTypeSingle[HLine] = tcell.RuneHLine
	borderTypeSingle[CloseLeft] = '['
	borderTypeSingle[CloseRight] = ']'
	borderTypeSingle[Close] = '■'
	borderTypeSingle[URCorner] = tcell.RuneURCorner
	borderTypeSingle[CaptionSpace] = ' '
	borderTypeSingle[LLCorner] = tcell.RuneLLCorner
	borderTypeSingle[LRCorner] = tcell.RuneLRCorner
	borderTypeSingle[VLine] = tcell.RuneVLine

	borderTypeDouble := make([]rune, 10)

	borderTypeDouble[ULCorner] = '╔'
	borderTypeDouble[HLine] = '═'
	borderTypeDouble[CloseLeft] = '['
	borderTypeDouble[CloseRight] = ']'
	borderTypeDouble[Close] = '■'
	borderTypeDouble[URCorner] = '╗'
	borderTypeDouble[CaptionSpace] = ' '
	borderTypeDouble[LLCorner] = '╚'
	borderTypeDouble[LRCorner] = '╝'
	borderTypeDouble[VLine] = '║'

	borderTypeEmpty := make([]rune, 10)

	borderTypeEmpty[ULCorner] = ' '
	borderTypeEmpty[HLine] = ' '
	borderTypeEmpty[CloseLeft] = '['
	borderTypeEmpty[CloseRight] = ']'
	borderTypeEmpty[Close] = '■'
	borderTypeEmpty[URCorner] = ' '
	borderTypeEmpty[CaptionSpace] = ' '
	borderTypeEmpty[LLCorner] = ' '
	borderTypeEmpty[LRCorner] = ' '
	borderTypeEmpty[VLine] = ' '

	bordersChars[BorderTypeSingle] = borderTypeSingle
	bordersChars[BorderTypeDouble] = borderTypeDouble
	bordersChars[BorderTypeEmpty] = borderTypeEmpty
}

// GetClientBounds return client bounds.
func (w Window) GetClientBounds() base.Rect {
	return calculateClientBounds(w.GetBounds(), w.Border.Type)
}

// Draw the view.
func (w Window) Draw() {
	if !w.GetVisible() {
		return
	}

	canvas := w.Canvas()

	canvas.SetBrush(tcell.StyleDefault.
		Foreground(w.GetForegroundColor()).
		Background(w.GetBackgroundColor()))

	// Draw background of window
	bounds := w.GetBounds()
	bounds.X = 0
	bounds.Y = 0

	canvas.Fill(bounds)

	drawTitle(&w)

	drawBottom(&w)

	drawBorderLeft(&w)

	drawBorderRight(&w)
}

func drawTitle(canvas base.TCanvas, w *Window) {
	titleBounds := base.Rect{
		X:      0,
		Y:      0,
		Width:  w.GetBounds().Width,
		Height: 1,
	}

	borderStyle := tcell.StyleDefault.
		Foreground(w.Border.ForegroundColor).
		Background(w.Border.BackgroundColor)

	canvas.SetBrush(borderStyle)

	DefaultDrawTitleBar(canvas, titleBounds, w.Caption, bordersChars[w.Border.Type])
}

func drawBottom(canvas base.TCanvas, w *Window) {
	bounds := w.GetBounds()

	bottomBorderBounds := base.Rect{
		X:      0,
		Y:      bounds.Height - 1,
		Width:  bounds.Width,
		Height: 1,
	}

	borderStyle := tcell.StyleDefault.
		Foreground(w.Border.ForegroundColor).
		Background(w.Border.BackgroundColor)

	canvas.SetBrush(borderStyle)

	DefaultDrawBottomBar(canvas, bottomBorderBounds, bordersChars[w.Border.Type])
}

func drawBorderLeft(absoluteBounds base.Rect, drawBounds base.Rect, w *Window) {
	leftBorderBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y + 1,
		Width:  1,
		Height: absoluteBounds.Height - 2,
	}

	borderStyle := tcell.StyleDefault.
		Foreground(w.Border.ForegroundColor).
		Background(w.Border.BackgroundColor)

	DefaultDrawLeftOrRightBorder(w.AppConfig().Screen, leftBorderBounds, drawBounds, w.Caption,
		borderStyle, bordersChars[w.Border.Type])
}

func drawBorderRight(absoluteBounds base.Rect, drawBounds base.Rect, w *Window) {
	leftBorderBounds := base.Rect{
		X:      absoluteBounds.X + absoluteBounds.Width - 1,
		Y:      absoluteBounds.Y + 1,
		Width:  1,
		Height: absoluteBounds.Height - 2,
	}

	borderStyle := tcell.StyleDefault.
		Foreground(w.Border.ForegroundColor).
		Background(w.Border.BackgroundColor)

	DefaultDrawLeftOrRightBorder(w.AppConfig().Screen, leftBorderBounds, drawBounds, w.Caption,
		borderStyle, bordersChars[w.Border.Type])
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
			for _, child := range w.Children() {
				if child.HandleMessage(base.BuildDrawMessage(child.Handler())) {
				}
			}
		}
	case base.WmChangeBounds:
		// Minimum Width/Height -> 2
		bounds := msg.Value.(base.Rect)

		bounds.Width = base.MaxInt(bounds.Width, 2)
		bounds.Height = base.MaxInt(bounds.Height, 2)

		w.SetBounds(bounds)
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

		// Redraw children.
		for _, child := range w.Children() {
			if child.HandleMessage(msg) {
			}
		}
	}

	return false
}

func calculateClientBounds(bounds base.Rect, borderType BorderType) base.Rect {
	/*
		  TODO window type
			switch borderType {
			case BorderTypeNone:
				// Remove titlebar
				bounds.Y++
				bounds.Height--
			default:
				// Remove titlebar and border
				bounds.Y++
				bounds.X++
				bounds.Height -= 2
				bounds.Width -= 2
			}*/

	bounds.Y = 1
	bounds.X = 1
	bounds.Height -= 2
	bounds.Width -= 2

	return bounds
}

// NewWindow create new window.
func NewWindow(name string, config base.ApplicationConfig, parentCanvas base.TCanvas) Window {
	w := Window{
		View:    base.NewView(name, config, parentCanvas),
		Caption: name,
		Border: WindowBorder{
			Type:            BorderTypeSingle,
			BackgroundColor: tcell.ColorGray,
			ForegroundColor: tcell.ColorWhite,
		},
	}

	w.SetBackgroundColor(tcell.ColorGray)
	w.SetForegroundColor(tcell.ColorWhite)

	return w
}

//------------------------------------------------------------------------------
// Helpher.

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

//------------------------------------------------------------------------------
// Default draw functions.

// DefaultDrawTitleBar default draw for title bar.
// Give ┌─[■]─ My title ─┐
func DefaultDrawTitleBar(canvas base.TCanvas, titleBounds base.Rect, caption string, borders []rune) {
	indexTitleBar := 0

	canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[ULCorner]) // [0]

	indexTitleBar++

	const minimumTitleBar = 7

	// Draw close only if available space for
	// ┌─[■]─┐
	// Need space before and after caption -> +2
	if titleBounds.Width >= minimumTitleBar {
		const minimumTitleBarWithCaption = minimumTitleBar + 2

		canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine]) // [1]
		indexTitleBar++
		canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CloseLeft]) // [2]
		indexTitleBar++
		// TODO use a button
		canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[Close]) // [3]
		indexTitleBar++
		canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CloseRight]) // [4]
		indexTitleBar++

		paddingLen := titleBounds.Width - len(caption) - 2
		paddingLenLeft := (paddingLen / 2)

		// Need space before and after caption -> +2
		// ─── Title ────
		if titleBounds.Width > minimumTitleBarWithCaption {
			// Border before caption
			// ──────── Hello
			for ; indexTitleBar < paddingLenLeft; indexTitleBar++ {
				canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine])
			}

			canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CaptionSpace])
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
				canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, c[i])
				indexTitleBar++
			}

			canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CaptionSpace])
			indexTitleBar++

			// Border after caption
			// Hello ────
			// -1 cause last char is corner
			for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
				canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine])
			}

			// TODO add up/down button

		} else {
			// No space to draw caption
			for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
				canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine])
			}
		}
	} else {
		// No space to draw close button
		for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
			canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine])
		}
	}

	canvas.PrintChar(titleBounds.X+indexTitleBar, titleBounds.Y, borders[URCorner])
}

// DefaultDrawBottomBar draw bottom border of window.
// Give └───────────┘
func DefaultDrawBottomBar(canvas base.Canvas, bottomBounds base.Rect, borders []rune) {
	indexTitleBar := 0

	canvas.PrintChar(bottomBounds.X+indexTitleBar, bottomBounds.Y, borders[LLCorner]) // [0]
	indexTitleBar++

	// -1 cause last char is corner
	for ; indexTitleBar < bottomBounds.Width-1; indexTitleBar++ {
		canvas.PrintChar(bottomBounds.X+indexTitleBar, bottomBounds.Y, borders[HLine])
	}

	canvas.PrintChar(bottomBounds.X+indexTitleBar, bottomBounds.Y, borders[LRCorner])
}

// DefaultDrawLeftOrRightBorder draw left or right border of window.
func DefaultDrawLeftOrRightBorder(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borders []rune) {
	titleBarLen := titleBounds.Height

	for indexTitleBar := 0; indexTitleBar < titleBarLen; indexTitleBar++ {
		screen.SetContent(titleBounds.X, titleBounds.Y+indexTitleBar, borders[VLine], nil, style)
	}
}
