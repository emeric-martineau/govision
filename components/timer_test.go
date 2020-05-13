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
	"testing"
	"time"

	base "github.com/emeric-martineau/govision"
	"github.com/gdamore/tcell"
)

func CreateTestApplicationConfig() base.ApplicationConfig {
	screen := tcell.NewSimulationScreen("")

	// Screen application.
	screenStyle := base.ApplicationStyle{
		Style:           tcell.StyleDefault,
		ForegroundColor: tcell.ColorWhite,
		BackgroundColor: tcell.ColorBlack,
	}

	return base.ApplicationConfig{
		ScreenStyle: screenStyle,
		Screen:      screen,
		Message:     base.NewBus(),
	}
}

func TestTime_Codecoverage(t *testing.T) {
	timer1 := NewTimer("timer1", 10*time.Millisecond, CreateTestApplicationConfig().Message)
	timer1.GetIntervale()
	timer1.SetEnabled(true)
	timer1.SetEnabled(false)
}

func TestTime_OnTimer(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	timer1 := NewTimer("timer1", 10*time.Millisecond, appConfig.Message)
	timer1.OnTimer = func(t *Timer) {
		appConfig.Message.Send(base.Message{})
	}

	timeOut := time.NewTimer(50 * time.Millisecond)

	timer1.SetEnabled(true)

	select {
	case <-timeOut.C:
		t.Error("OnTimer not called!")
	case <-*appConfig.Message.Channel():
	}
}

func TestTime_WmEnable(t *testing.T) {
	isCalled := false
	appConfig := CreateTestApplicationConfig()

	app := base.NewApplication(appConfig)

	c1 := base.NewView("c1", appConfig.Message, app.Canvas())
	c1.SetOnReceiveMessage(func(c base.TComponent, m base.Message) bool {
		switch m.Type {
		case base.WmTimer:
			isCalled = true
			appConfig.Message.Send(base.Message{
				Handler: base.ApplicationHandler(),
				Type:    base.WmQuit,
			})
		}

		return false
	})

	// AddChild is not necessary
	appConfig.Message.Send(base.Message{
		Handler: base.ApplicationHandler(),
		Type:    base.WmCreate,
		Value:   &c1,
	})

	timer1 := NewTimer("timer1", 0, appConfig.Message)
	timer1.SetParent(&c1)

	app.AddWindow(&c1)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		timer1.SetIntervale(10 * time.Millisecond)
		timer1.SetEnabled(true)

		app.Run()

		if !isCalled {
			t.Error("OnTimer not called!")
		}
	}
}

func TestTime_SetInterval_if_enable(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	timer1 := NewTimer("timer1", 5*time.Millisecond, appConfig.Message)

	OnTimer2 := func(t *Timer) {
		appConfig.Message.Send(base.Message{})
	}

	OnTimer1 := func(t *Timer) {
		timer1.OnTimer = OnTimer2
		t.SetIntervale(7 * time.Millisecond)
	}

	timer1.OnTimer = OnTimer1

	timeOut := time.NewTimer(50 * time.Millisecond)

	timer1.SetEnabled(true)

	select {
	case <-timeOut.C:
		t.Error("OnTimer not called!")
	case <-*appConfig.Message.Channel():
	}
}

func TestTime_OnTimer_disable_timer(t *testing.T) {
	appConfig := CreateTestApplicationConfig()

	timer1 := NewTimer("timer1", 10*time.Millisecond, appConfig.Message)
	timer1.OnTimer = func(t *Timer) {
		t.SetEnabled(false)
	}

	timeOut := time.NewTimer(50 * time.Millisecond)

	timer1.SetEnabled(true)

	select {
	case <-timeOut.C:
	case <-*appConfig.Message.Channel():
		t.Error("OnTimer not called!")
	}
}
