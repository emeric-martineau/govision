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
	"time"

	base "govision"
)

// OnTimer is callback when timer is done.
type OnTimer func(*Timer)

// Timer is the base object of all widget.
type Timer struct {
	// Interval timer
	interval time.Duration
	// Use to cancel timer
	canceled chan bool
	// OnTimer is callback when timer is done.
	OnTimer OnTimer

	base.Component
}

// GetIntervale return interval value.
func (t *Timer) GetIntervale() time.Duration {
	return t.interval
}

// SetIntervale set new interval and reset timer.
func (t *Timer) SetIntervale(interval time.Duration) {
	// Cancel time
	t.canceled <- true
	t.interval = interval

	go run(t)
}

// TODO WmTimer send to parent if OnTimer is nil

// SetEnabled active or disable timer.
func (t *Timer) SetEnabled(status bool) {
	if status && !t.GetEnabled() {
		// Start timer
		go run(t)
	} else if !status && t.GetEnabled() {
		// Disable timer
		t.canceled <- true
	}

	t.Component.SetEnabled(status)
}

func run(t *Timer) {
	var timer *time.Timer

loop:
	for {
		timer = time.NewTimer(t.interval)

		select {
		case <-timer.C:
			if t.OnTimer != nil {
				t.OnTimer(t)
			}

			if !t.GetEnabled() {
				break loop
			}
		case <-t.canceled:
			break loop
		}
	}
}

// NewTimer create new timer.
func NewTimer(name string, interval time.Duration) Timer {
	return Timer{
		Component: base.NewComponent(name),
		interval:  interval,
		canceled:  make(chan bool),
	}
}
