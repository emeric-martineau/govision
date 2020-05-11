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
	"github.com/google/uuid"
)

// OnDraw is call when change enable.
type OnDraw func(TView)

// View is the base object of all visual widget.
type View struct {
	component       Component
	canvas          TCanvas
	bounds          Rect
	focused         bool
	visible         bool
	backgroundColor tcell.Color
	foregroundColor tcell.Color
	// To overide draw for custom draw for example.
	onDraw OnDraw
}

// Name of component.
func (v *View) Name() string {
	return v.component.Name()
}

// Handler of component (UUID).
func (v *View) Handler() uuid.UUID {
	return v.component.Handler()
}

// SetEnabled enable component.
func (v *View) SetEnabled(e bool) {
	v.component.SetEnabled(e)
}

// GetEnabled return is component is enable.
func (v *View) GetEnabled() bool {
	return v.component.GetEnabled()
}

// SetParent set parent component.
func (v *View) SetParent(p TComponent) {
	v.component.SetParent(p)
}

// GetParent get parent component.
func (v *View) GetParent() TComponent {
	return v.component.GetParent()
}

// AddChild add a child to component.
func (v *View) AddChild(c TComponent) {
	v.component.AddChild(c)
}

// RemoveChild remove a child to component.
func (v *View) RemoveChild(c TComponent) {
	v.component.RemoveChild(c)
}

// Children return list of children of component.
func (v *View) Children() []TComponent {
	return v.component.Children()
}

// SetZorder set the new odrer of draw.
func (v *View) SetZorder(o int) {
	v.component.SetZorder(o)
}

// GetZorder return odrer of draw.
func (v *View) GetZorder() int {
	return v.component.GetZorder()
}

// SetOnDraw set ondraw callback.
func (v *View) SetOnDraw(f OnDraw) {
	v.onDraw = f
}

// GetOnDraw get ondraw callback.
func (v *View) GetOnDraw() OnDraw {
	return v.onDraw
}

// SetBounds set view size.
func (v *View) SetBounds(r Rect) {
	v.bounds = r

	v.canvas.UpdateBounds(r)
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

// GetMessageBus return application config.
func (v *View) GetMessageBus() Bus {
	return v.component.GetMessageBus()
}

// Draw the view.
func (v *View) Draw() {
	if !v.visible {
		return
	}

	v.canvas.SetBrush(tcell.StyleDefault.
		Foreground(v.foregroundColor).
		Background(v.backgroundColor))

	bounds := v.GetBounds()
	bounds.X = 0
	bounds.Y = 0

	v.canvas.Fill(bounds)
}

// Canvas of view.
func (v *View) Canvas() TCanvas {
	return v.canvas
}

// ClientCanvas return client canvas of view.
func (v *View) ClientCanvas() TCanvas {
	return v.canvas.CreateCanvasFrom(v.GetClientBounds())
}

//------------------------------------------------------------------------------
// Internal function.

// Manage message if it's for me.
// Return true to stop message propagation.
func (v *View) manageMyMessage(msg Message) {
	switch msg.Type {
	case WmDraw:
		if v.onDraw != nil {
			v.onDraw(v)
		} else {
			v.Draw()
			// Redraw children.
			v.component.HandleMessage(BuildDrawMessage(BroadcastHandler()))
		}
	case WmChangeBounds:
		v.SetBounds(msg.Value.(Rect))
		// Redraw all components cause maybe overide a component with Zorder
		v.component.message.Send(BuildDrawMessage(BroadcastHandler()))
	default:
		v.component.HandleMessage(msg)
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
	return v.component.HandleMessage(msg)
}

//------------------------------------------------------------------------------
// Constrcutor.

// NewView create new timer.
func NewView(name string, message Bus, parentCanvas TCanvas) View {
	return View{
		component: NewComponent(name, message),
		canvas:    parentCanvas.CreateCanvasFrom(Rect{0, 0, 0, 0}),
	}
}
