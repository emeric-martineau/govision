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
const WmNull uint = 0x00

// WmEnable send to enable or disable a component.
const WmEnable uint = 0x01

// WmKey when a key is pressed.
const WmKey uint = 0x02

// WmScreenResize when a screen size change.
const WmScreenResize uint = 0x03

// WmDraw draw component. Becarefull! If you send to one component, only this
// component and his child draw.
const WmDraw uint = 0x04

// WmZorderChange send to parent when you change Zorder of children.
const WmZorderChange uint = 0x05

// WmQuit force application to shutdown.
const WmQuit uint = 0x06

// WmChangeBounds send to component to change size, or move.
const WmChangeBounds uint = 0x07

// WmUser allow user to have own message.
const WmUser uint = ^uint(0) / 2

// BuildKeyMessage build a message for keyboard event.
func BuildKeyMessage(event *tcell.EventKey) Message {
	return Message{
		Handler: BroadcastHandler(),
		Type:    WmKey,
		Value:   event,
	}
}

// BuildScreenResizeMessage build a resize message broadcast.
func BuildScreenResizeMessage(s tcell.Screen) Message {
	width, height := s.Size()

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

// BuildChangeBounds return a message to change bounds of component.
func BuildChangeBounds(handler uuid.UUID, bounds Rect) Message {
	return Message{
		Handler: handler,
		Type:    WmChangeBounds,
		Value:   bounds,
	}
}
