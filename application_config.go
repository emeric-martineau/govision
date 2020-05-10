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
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

// ApplicationConfig is configuration of application.
type ApplicationConfig struct {
	// Default style screen application.
	ScreenStyle ApplicationStyle // TODO remove cause use canvas
	// Screen of application.
	Screen tcell.Screen // TODO remove cause use canvas
	// Message bus.
	Message Bus
}

// ApplicationStyle is style of screen for application.
type ApplicationStyle struct {
	// Default screen style.
	Style tcell.Style
	// Default screen text color.
	ForegroundColor tcell.Color
	// Default screen background color.
	BackgroundColor tcell.Color
}

// CreateDefaultApplicationConfig create application config for almost case.
func CreateDefaultApplicationConfig() ApplicationConfig {
	screen, e := tcell.NewScreen()

	if e != nil {
		fmt.Fprintf(os.Stderr, "Error when create screen of application: %+v\n", e)
		os.Exit(1)
	}

	// Screen application.
	screenStyle := ApplicationStyle{
		Style:           tcell.StyleDefault,
		ForegroundColor: tcell.ColorWhite,
		BackgroundColor: tcell.ColorBlack,
	}

	return ApplicationConfig{
		ScreenStyle: screenStyle,
		Screen:      screen,
		Message:     NewBus(),
	}
}
