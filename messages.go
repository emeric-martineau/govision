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

// Message is data structure for any message.
type Message struct {
	Handler uuid.UUID
	Type    uint
	Value   interface{}
}

// WmNull is empty message, ignore it (internal use only). This is never use by
// GoVision.
const WmNull uint = 0

// WmEnable send to enable or disable a component.
const WmEnable uint = 1

// WmKey when a key is pressed.
const WmKey uint = 2

// WmScreenResize when a screen size change.
const WmScreenResize uint = 3

// WmDraw draw component. Becarefull! If you send to one component, only this
// component and his child draw.
const WmDraw uint = 4

// WmZorderChange send to parent when you change Zorder of children.
const WmZorderChange uint = 5

// WmQuit force application to shutdown.
const WmQuit uint = 6

// WmChangeBounds send to component to change size, or move.
const WmChangeBounds uint = 7

// WmTimer send to Timer parent component if OnTimer is nil.
const WmTimer uint = 8

// WmCreate sent when you have create windows and want add in list.
const WmCreate uint = 9

// WmDestroy sent when you want remove windows from list and let GC remove it.
const WmDestroy uint = 10

// WmMouse send when mouse occure. Generally manage by Application struct.
// If nothing can be done with, send to Window struct.
const WmMouse uint = 11

// WmLButtonDown send when left button pressed on window beneath cursor.
const WmLButtonDown uint = 12

// WmLButtonUp send when left button released on window beneath cursor.
const WmLButtonUp uint = 13

// WmRButtonDown send when right button pressed on window beneath cursor.
const WmRButtonDown uint = 14

// WmRButtonUp send when right button released on window beneath cursor.
const WmRButtonUp uint = 15

// WmActivate Sent to both the window being activated and the window being deactivated.
// Value can be WaActive or WaInactive.
const WmActivate uint = 16

// WmMouseEnter sent when mouse enter to TView.
const WmMouseEnter uint = 17

// WmMouseLeave sent when mouse leave to TView.
const WmMouseLeave uint = 18

// WmUser allow user to have own message.
const WmUser uint = ^uint(0) / 2

// WaInactive Deactivated.
const WaInactive uint = 0

// WaActive Activated.
const WaActive uint = 1

// BuildKeyMessage build a message for keyboard event.
func BuildKeyMessage(event *tcell.EventKey) Message {
	return Message{
		Handler: BroadcastHandler(),
		Type:    WmKey,
		Value:   event,
	}
}

// BuildScreenResizeMessage build a resize message broadcast.
func BuildScreenResizeMessage(screen tcell.Screen) Message {
	width, height := screen.Size()

	r := Rect{
		X:      0,
		Y:      0,
		Width:  width,
		Height: height,
	}

	return Message{
		Handler: BroadcastHandler(),
		Type:    WmScreenResize,
		Value:   r,
	}
}

// BuildEmptyMessage build empty message.
func BuildEmptyMessage() Message {
	return Message{
		Handler: BroadcastHandler(),
		Type:    WmNull,
	}
}

// BuildDrawMessage return a message to force a component to redraw.
func BuildDrawMessage(handler uuid.UUID) Message {
	return Message{
		Handler: handler,
		Type:    WmDraw,
	}
}

// BuildZorderMessage return a message to change Zorder.
func BuildZorderMessage(handler uuid.UUID) Message {
	return Message{
		Handler: handler,
		Type:    WmZorderChange,
	}
}

// BuildChangeBoundsMessage return a message to change bounds of component.
func BuildChangeBoundsMessage(handler uuid.UUID, bounds Rect) Message {
	return Message{
		Handler: handler,
		Type:    WmChangeBounds,
		Value:   bounds,
	}
}

// BuildActivateMessage send message to windows gains focus.
func BuildActivateMessage(handler uuid.UUID) Message {
	return Message{
		Handler: handler,
		Type:    WmActivate,
		Value:   WaActive,
	}
}

// BuildDesactivateMessage send message to windows gains focus.
func BuildDesactivateMessage(handler uuid.UUID) Message {
	return Message{
		Handler: handler,
		Type:    WmActivate,
		Value:   WaInactive,
	}
}

// BuildClickMouseMessage send message to windows gains focus.
func BuildClickMouseMessage(handler uuid.UUID, ev *tcell.EventMouse, side uint) Message {
	return Message{
		Handler: handler,
		Type:    side,
		Value:   ev,
	}
}

// BuildMouseEnterMessage send message to windows when mouse enter.
func BuildMouseEnterMessage(handler uuid.UUID, x, y int) Message {
	return Message{
		Handler: handler,
		Type:    WmMouseEnter,
		Value: Rect{
			X: x,
			Y: y,
		},
	}
}

// BuildMouseLeaveMessage send message to windows when mouse enter.
func BuildMouseLeaveMessage(handler uuid.UUID) Message {
	return Message{
		Handler: handler,
		Type:    WmMouseLeave,
	}
}
