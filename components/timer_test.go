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

func TestTime_OnTimer_function(t *testing.T) {
	isCalled := false

	timer1 := NewTimer("timer1", 10*time.Millisecond)
	timer1.OnTimer = func(t *Timer) {
		isCalled = true
		base.SendMessage(base.Message{
			Handler: base.BroadcastHandler(),
			Type:    base.WmQuit,
		})
	}

	app, _ := base.NewApplication(&timer1)

	if e := app.Init(); e != nil {
		t.Error("Cannot initialize screen")
	} else {
		timer1.SetEnabled(true)

		app.Run()

		if !isCalled {
			t.Error("OnTimer not called!")
		}
	}
}
