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
)

func TestTime_Codecoverage(t *testing.T) {
	timer1 := NewTimer("timer1", 10*time.Millisecond)
	timer1.GetIntervale()
	timer1.SetEnabled(true)
	timer1.SetEnabled(false)
}

func TestTime_OnTimer(t *testing.T) {
	busChannel := make(chan bool)

	timer1 := NewTimer("timer1", 10*time.Millisecond)
	timer1.OnTimer = func(t *Timer) {
		busChannel <- true
	}

	timeOut := time.NewTimer(50 * time.Millisecond)

	timer1.SetEnabled(true)

	select {
	case <-timeOut.C:
		t.Error("OnTimer not called!")
	case <-busChannel:
	}
}

func TestTime_WmEnable(t *testing.T) {
	isCalled := false

	c1 := base.NewComponent("c1")
	c1.OnReceiveMessage = func(c base.TComponent, m base.Message) bool {
		switch m.Type {
		case base.WmTimer:
			isCalled = true
			base.SendMessage(base.Message{
				Handler: base.ApplicationHandler(),
				Type:    base.WmQuit,
			})
		}

		return false
	}

	timer1 := NewTimer("timer1", 0)
	timer1.SetParent(&c1)

	// AddChild is not necessary
	base.SendMessage(base.Message{
		Handler: base.ApplicationHandler(),
		Type:    base.WmCreate,
		Value:   &c1,
	})

	app, _ := base.NewApplication(&c1)

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
	busChannel := make(chan bool)

	timer1 := NewTimer("timer1", 5*time.Millisecond)

	OnTimer2 := func(t *Timer) {
		busChannel <- true
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
	case <-busChannel:
	}
}

func TestTime_OnTimer_disable_timer(t *testing.T) {
	busChannel := make(chan bool)

	timer1 := NewTimer("timer1", 10*time.Millisecond)
	timer1.OnTimer = func(t *Timer) {
		t.SetEnabled(false)
	}

	timeOut := time.NewTimer(50 * time.Millisecond)

	timer1.SetEnabled(true)

	select {
	case <-timeOut.C:
	case <-busChannel:
		t.Error("OnTimer not called!")
	}
}
