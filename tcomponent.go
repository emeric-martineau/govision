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
	"github.com/google/uuid"
)

// TComponent is the base object of all widget.
type TComponent interface {
	// Call when receive a message.
	// Return true if stop loop because message are manage by client.
	HandleMessage(Message) bool
	// Name of component.
	Name() string
	// Handler of component (UUID).
	Handler() uuid.UUID
	// Enable component.
	SetEnabled(bool)
	GetEnabled() bool
	// Parent component.
	SetParent(TComponent)
	GetParent() TComponent
	// Children component.
	AddChild(TComponent)
	RemoveChild(TComponent)
	Children() []TComponent
	// Odrer of draw.
	SetZorder(int)
	GetZorder() int
}
