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

import "testing"

// Test if function works :)
func TestComponent_Dummy_for_code_coverage(t *testing.T) {
	c := NewComponent("You know my name")
	c.Name()
	c.Handler()
	c.GetEnabled()
	c.SetZorder(1)
	c.GetZorder()
}

// Test if function works :)
func TestComponent_Enable_without_OnEnabled(t *testing.T) {
	c := NewComponent("You know my name")
	c.SetEnabled(true) // By default it's false

	if c.GetEnabled() != true {
		t.Errorf("Cannot set enable component")
	}
}

func TestComponent_Enable_with_OnEnabled(t *testing.T) {
	c := NewComponent("You know my name")
	c.OnEnabled = func(c TComponent, status bool) bool {
		return false
	}

	c.SetEnabled(true) // By default it's false

	if c.GetEnabled() == true {
		t.Errorf("OnEnable() not called!")
	}
}

func TestComponent_AddChild(t *testing.T) {
	c1 := NewComponent("c1")
	c2 := NewComponent("c2")

	c2.AddChild(&c1)

	children := c2.Children()

	if len(children) != 1 {
		t.Errorf("Normally, only one child, found %d", len(children))
	}

	if children[0] != &c1 {
		t.Errorf("Normally, 'c1' but found '%s'", children[0].Name())
	}
}

func TestComponent_RemoveChild(t *testing.T) {
	c1 := NewComponent("c1")
	c2 := NewComponent("c2")
	c3 := NewComponent("c3")
	c4 := NewComponent("c4")

	c2.AddChild(&c1)
	c2.AddChild(&c3)
	c2.AddChild(&c4)

	c2.RemoveChild(&c3)

	children := c2.Children()

	if len(children) != 2 {
		t.Errorf("Normally, only two children, found %d", len(children))
	}

	if children[0] != &c1 {
		t.Errorf("Normally, 'c1' but found '%s'", children[0].Name())
	}

	if children[1] != &c4 {
		t.Errorf("Normally, 'c4' but found '%s'", children[1].Name())
	}
}

func TestComponent_RemoveChild_not_found_component(t *testing.T) {
	c1 := NewComponent("c1")
	c2 := NewComponent("c2")
	c3 := NewComponent("c3")
	c4 := NewComponent("c4")

	c2.AddChild(&c1)
	//c2.AddChild(&c3)
	c2.AddChild(&c4)

	children := c2.Children()

	c2.RemoveChild(&c3)

	if len(children) != 2 {
		t.Errorf("Normally, only two children, found %d", len(children))
	}

	if children[0] != &c1 {
		t.Errorf("Normally, 'c1' but found '%s'", children[0].Name())
	}

	if children[1] != &c4 {
		t.Errorf("Normally, 'c4' but found '%s'", children[1].Name())
	}
}

func TestComponent_Message_draw_send_to_a_child(t *testing.T) {
	isCalled := false

	// Use OnReceiveMessage on child
	c1 := NewComponent("c1")
	c2 := NewComponent("c2")

	c1.AddChild(&c2)

	c2.OnReceiveMessage = func(c TComponent, m Message) bool {
		isCalled = true

		if c != &c2 {
			t.Errorf("Normally, 'c1' but found '%s'", c.Name())
		}

		if m.Type != WmDraw {
			t.Errorf("Normally, 'WmDraw' but found '%d'", m.Type)
		}

		return false
	}

	c1.HandleMessage(Message{
		Handler: c1.Handler(),
		Type:    WmDraw,
	})

	if !isCalled {
		t.Errorf("Message not send to child")
	}
}

func TestComponent_Message_draw_broadcast(t *testing.T) {
	// Use OnReceiveMessage on main component and child
	isCalledMain := false
	isCalledChild := false

	c1 := NewComponent("c1")
	c2 := NewComponent("c2")

	c1.AddChild(&c2)

	c1.OnReceiveMessage = func(c TComponent, m Message) bool {
		isCalledMain = true

		if c != &c1 {
			t.Errorf("Normally, 'c1' but found '%s'", c.Name())
		}

		if m.Type != WmDraw {
			t.Errorf("Normally, 'WmDraw' but found '%d'", m.Type)
		}

		return false
	}

	c2.OnReceiveMessage = func(c TComponent, m Message) bool {
		isCalledChild = true

		if c != &c2 {
			t.Errorf("Normally, 'c2' but found '%s'", c.Name())
		}

		if m.Type != WmDraw {
			t.Errorf("Normally, 'WmDraw' but found '%d'", m.Type)
		}

		return false
	}

	c1.HandleMessage(Message{
		Handler: BroadcastHandler(),
		Type:    WmDraw,
	})

	if !isCalledMain {
		t.Error("Message not send to main")
	}

	if !isCalledChild {
		t.Error("Message not send to child")
	}
}

func TestComponent_Message_enable(t *testing.T) {
	// Use OnReceiveMessage on main component
	c1 := NewComponent("c1")
	c1.SetEnabled(true) // By default it's false

	c1.OnEnabled = func(c TComponent, status bool) bool {
		if status == false {
			t.Error("Bad value called!")
		}

		return false
	}

	c1.HandleMessage(Message{
		Handler: c1.Handler(),
		Type:    WmEnable,
		Value:   true,
	})

	if c1.GetEnabled() == true {
		t.Error("OnEnable() not called!")
	}
}

func TestComponent_Message_zorder(t *testing.T) {
	c := NewComponent("main")

	c1 := NewComponent("1")
	c1.SetZorder(2)
	c2 := NewComponent("2")
	c2.SetZorder(1)
	c3 := NewComponent("3")
	c3.SetZorder(3)

	c.AddChild(&c1)
	c.AddChild(&c2)
	c.AddChild(&c3)

	c.HandleMessage(Message{
		Handler: c1.Handler(),
		Type:    WmZorderChange,
	})

	children := c.Children()

	compare(children[0], &c2, t)
	compare(children[1], &c1, t)
	compare(children[2], &c3, t)
}
