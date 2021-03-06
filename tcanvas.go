package base

import "github.com/gdamore/tcell"

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

// TCanvas is virtal screen to draw components.
type TCanvas interface {
	// Set brush property.
	SetBrush(b tcell.Style)
	// Create a canvas.
	CreateCanvasFrom(r Rect) TCanvas
	// Print char with current brush property.
	PrintChar(x int, y int, char rune)
	// Print char with brush.
	PrintCharWithBrush(x int, y int, char rune, brush tcell.Style)
	// Update position.
	UpdateBounds(r Rect)
	// Fill zone of canvas.
	Fill(bounds Rect)
}
