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
	"strings"

	"github.com/gdamore/tcell"
)

// Internal screen.
var screen tcell.Screen
var appScreen ApplicationScreen

// Inits a new screen.
func init() {
	if strings.HasSuffix(os.Args[0], ".test") {
		screen = tcell.NewSimulationScreen("")
	} else {
		var e error
		screen, e = tcell.NewScreen()

		if e != nil {
			fmt.Fprintf(os.Stderr, "Error when create screen of application: %+v\n", e)
			os.Exit(1)
		}
	}

	// Screen application.
	appScreen = ApplicationScreen{
		Style:           tcell.StyleDefault,
		ForegroundColor: tcell.ColorWhite,
		BackgroundColor: tcell.ColorBlack,
	}
}

// ApplicationScreen is base struct for screen of application.
type ApplicationScreen struct {
	// Default screen style.
	Style tcell.Style
	// Default screen text color.
	ForegroundColor tcell.Color
	// Default screen background color.
	BackgroundColor tcell.Color
}

// Screen return screen to draw
func (a ApplicationScreen) Screen() tcell.Screen {
	return screen
}

// AppScreen return struct to draw on screen.
func AppScreen() *ApplicationScreen {
	return &appScreen
}
