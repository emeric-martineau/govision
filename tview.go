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

// TView is the base object of all visual widget.
type TView interface {
	TComponent
	// External component size.
	SetBounds(Rect)
	GetBounds() Rect
	// Focus component.
	SetFocused(bool)
	GetFocused() bool
	// Visible component.
	SetVisible(bool)
	GetVisible() bool
	// Client size.
	GetClientBounds() Rect
	// Color of component.
	GetBackgroundColor() tcell.Color
	SetBackgroundColor(tcell.Color)
	GetForegroundColor() tcell.Color
	SetForegroundColor(tcell.Color)
	// Draw component
	Draw()
	// Screen to draw.
	GetScreen() tcell.Screen
}
