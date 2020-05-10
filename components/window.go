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

	style := tcell.StyleDefault.
		Foreground(w.GetForegroundColor()).
		Background(w.GetBackgroundColor())

	// Get parent X and Y
	absoluteBounds := base.CalculateAbsolutePosition(&w)
	//base.PrintStringOnScreen(w.AppConfig().Screen, tcell.ColorBlack, tcell.ColorWhite, absoluteBounds.X-3, absoluteBounds.Y-3, fmt.Sprintf("%s: %+v", w.Name(), absoluteBounds))
	// Get real zone whe can draw
	drawBounds := base.CalculateDrawZone(&w)
	//base.PrintStringOnScreen(w.AppConfig().Screen, tcell.ColorBlack, tcell.ColorWhite, absoluteBounds.X-2, absoluteBounds.Y-2, fmt.Sprintf("%s: %+v", w.Name(), drawBounds))
	// Get client bould to draw
	clientBounds := calculateClientBounds(absoluteBounds, w.Border.Type)
	clientBounds.X += absoluteBounds.X
	clientBounds.Y += absoluteBounds.Y

	//base.PrintStringOnScreen(w.AppConfig().Screen, tcell.ColorBlack, tcell.ColorWhite, absoluteBounds.X-1, absoluteBounds.Y-1, fmt.Sprintf("%s: %+v", w.Name(), clientBounds))

	// Draw background of window
	base.Fill(w.AppConfig().Screen, clientBounds, drawBounds, style)

	drawTitle(absoluteBounds, drawBounds, &w)

	drawBottom(absoluteBounds, drawBounds, &w)

	drawBorderLeft(absoluteBounds, drawBounds, &w)

	drawBorderRight(absoluteBounds, drawBounds, &w)
}

func drawTitle(absoluteBounds base.Rect, drawBounds base.Rect, w *Window) {
	titleBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y,
		Width:  absoluteBounds.Width,
		Height: 1,
	}

	borderStyle := tcell.StyleDefault.
		Foreground(w.Border.ForegroundColor).
		Background(w.Border.BackgroundColor)

	DefaultDrawTitleBar(w.AppConfig().Screen, titleBounds, drawBounds, w.Caption,
		borderStyle, bordersChars[w.Border.Type])
}

func drawBottom(absoluteBounds base.Rect, drawBounds base.Rect, w *Window) {
	bottomBorderBounds := base.Rect{
		X:      absoluteBounds.X,
		Y:      absoluteBounds.Y + absoluteBounds.Height - 1,
		Width:  absoluteBounds.Width,
		Height: 1,
	}

	borderStyle := tcell.StyleDefault.
		Foreground(w.Border.ForegroundColor).
		Background(w.Border.BackgroundColor)

	DefaultDrawBottomBar(w.AppConfig().Screen, bottomBorderBounds, drawBounds, w.Caption,
		borderStyle, bordersChars[w.Border.Type])
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
func DefaultDrawTitleBar(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borders []rune) {
	indexTitleBar := 0

	base.PrintChar(screen, titleBounds.X+indexTitleBar, titleBounds.Y, borders[ULCorner], style, drawBounds) // [1]
	indexTitleBar++

	const minimumTitleBar = 7

	// Draw close only if available space for
	// ┌─[■]─┐
	// Need space before and after caption -> +2
	if titleBounds.Width >= minimumTitleBar {
		const minimumTitleBarWithCaption = minimumTitleBar + 2

		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine], nil, style) // [1]
		indexTitleBar++
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CloseLeft], nil, style) // [2]
		indexTitleBar++
		// TODO use a button
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[Close], nil, style) // [3]
		indexTitleBar++
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CloseRight], nil, style) // [4]
		indexTitleBar++

		paddingLen := titleBounds.Width - len(caption) - 2
		paddingLenLeft := (paddingLen / 2)

		// Need space before and after caption -> +2
		// ─── Title ────
		if titleBounds.Width > minimumTitleBarWithCaption {
			// Border before caption
			// ──────── Hello
			for ; indexTitleBar < paddingLenLeft; indexTitleBar++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine], nil, style)
			}

			screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CaptionSpace], nil, style)
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

			screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[CaptionSpace], nil, style)
			indexTitleBar++

			// Border after caption
			// Hello ────
			// -1 cause last char is corner
			for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine], nil, style)
			}

			// TODO add up/down button

		} else {
			// No space to draw caption
			for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
				screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine], nil, style)
			}
		}
	} else {
		// No space to draw close button
		for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
			screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine], nil, style)
		}
	}

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[URCorner], nil, style)
}

// DefaultDrawBottomBar draw bottom border of window.
// Give └───────────┘
func DefaultDrawBottomBar(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borders []rune) {
	indexTitleBar := 0

	// TODO if not fully visible

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[LLCorner], nil, style) // [0]
	indexTitleBar++

	// -1 cause last char is corner
	for ; indexTitleBar < titleBounds.Width-1; indexTitleBar++ {
		screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[HLine], nil, style)
	}

	screen.SetContent(titleBounds.X+indexTitleBar, titleBounds.Y, borders[LRCorner], nil, style)
}

// DefaultDrawLeftOrRightBorder draw left or right border of window.
func DefaultDrawLeftOrRightBorder(screen tcell.Screen, titleBounds base.Rect, drawBounds base.Rect, caption string, style tcell.Style, borders []rune) {
	titleBarLen := titleBounds.Height

	for indexTitleBar := 0; indexTitleBar < titleBarLen; indexTitleBar++ {
		screen.SetContent(titleBounds.X, titleBounds.Y+indexTitleBar, borders[VLine], nil, style)
	}
}
