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
	"github.com/gdamore/tcell"
)

// OnDraw is call when change enable.
type OnDraw func(TView)

// View is the base object of all visual widget.
type View struct {
	Component
	bounds          Rect
	focused         bool
	visible         bool
	backgroundColor tcell.Color
	foregroundColor tcell.Color
	// To overide draw for custom draw for example.
	OnDraw OnDraw
}

// SetBounds set view size.
func (v *View) SetBounds(r Rect) {
	v.bounds = r
}

// GetBounds return view size.
func (v *View) GetBounds() Rect {
	return v.bounds
}

// SetFocused of component.
func (v *View) SetFocused(f bool) {
	v.focused = f
}

// GetFocused return true if component has focus.
func (v *View) GetFocused() bool {
	return v.focused
}

// SetVisible if component is visible.
func (v *View) SetVisible(s bool) {
	v.visible = s
}

// GetVisible if component is visible.
func (v *View) GetVisible() bool {
	return v.visible
}

// GetClientBounds return client size.
func (v *View) GetClientBounds() Rect {
	return Rect{
		X:      0,
		Y:      0,
		Width:  v.bounds.Width,
		Height: v.bounds.Height,
	}
}

// GetBackgroundColor return background color.
func (v *View) GetBackgroundColor() tcell.Color {
	return v.backgroundColor
}

// SetBackgroundColor change background color.
func (v *View) SetBackgroundColor(c tcell.Color) {
	v.backgroundColor = c
}

// GetForegroundColor return text color.
func (v *View) GetForegroundColor() tcell.Color {
	return v.foregroundColor
}

// SetForegroundColor change text color.
func (v *View) SetForegroundColor(c tcell.Color) {
	v.foregroundColor = c
}

// Draw the view.
func (v *View) Draw() {
	if !v.visible {
		return
	}

	style := tcell.StyleDefault.
		Foreground(v.foregroundColor).
		Background(v.backgroundColor)

	// Get parent X and Y
	bounds := CalculateDrawZone(v)

	// If component is more biggest than parent
	startX := bounds.X
	endX := bounds.X + bounds.Width

	startY := bounds.Y
	endY := bounds.Y + bounds.Height

	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			v.AppConfig().Screen.
				SetContent(x, y, ' ', nil, style)
		}
	}
}

// Manage message if it's for me.
// Return true to stop message propagation.
func (v *View) manageMyMessage(msg Message) {
	switch msg.Type {
	case WmDraw:
		if v.OnDraw != nil {
			v.OnDraw(v)
		} else {
			v.Draw()
			// Redraw children.
			v.Component.HandleMessage(BuildDrawMessage(BroadcastHandler()))
		}
	case WmChangeBounds:
		v.bounds = msg.Value.(Rect)
		// Redraw all components cause maybe overide a component with Zorder
		v.AppConfig().Message.Send(BuildDrawMessage(BroadcastHandler()))
	}
}

// HandleMessage is use to manage message.
func (v *View) HandleMessage(msg Message) bool {
	switch msg.Handler {
	case v.Handler():
		v.manageMyMessage(msg)
		return true
	case BroadcastHandler():
		v.manageMyMessage(msg)
	}

	// Because Component send message to child if broadcast or draw.
	return v.Component.HandleMessage(msg)
}

// NewView create new timer.
func NewView(name string, config ApplicationConfig) View {
	return View{
		Component: NewComponent(name, config),
	}
}
