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
	"sort"

	"github.com/google/uuid"
)

// OnEnabled is call when change enable.
type OnEnabled func(TComponent, bool) bool

// OnReceiveMessage is call when component receive message and you want overide
// behavior.
// Return true to stop message propagation.
type OnReceiveMessage func(TComponent, Message) bool

// Component is the base object of all widget.
type Component struct {
	// Name for debugging for example.
	name string
	// Unique identifier
	handler uuid.UUID
	// If component is enable.
	enabled bool
	// Parent of component.
	parent TComponent
	// List of children component. List is always order by zorder.
	children []TComponent
	// If reveive message from HandlerMessage.
	OnReceiveMessage OnReceiveMessage
	// Call when component is enable.
	OnEnabled OnEnabled
	// Order to display.
	zorder int
}

// Manage message if it's for me.
// Return true to stop message propagation.
func (c *Component) manageMyMessage(msg Message) {
	if c.OnReceiveMessage != nil {
		c.OnReceiveMessage(c, msg)
	} else {
		switch msg.Type {
		case WmZorderChange:
			c.reorderChildren()
		case WmDraw:
			c.drawChildren()
		case WmEnable:
			c.SetEnabled(msg.Value.(bool))
		}
	}
}

// HandleMessage is use to manage message and give message to children.
func (c *Component) HandleMessage(msg Message) bool {
	var isStop bool

	switch msg.Handler {
	case c.handler:
		// For me
                c.manageMyMessage(msg)
		isStop = true
	case BroadcastHandler():
		// For me and my children
		isStop = false
		c.manageMyMessage(msg)

		// For my children
		for _, child := range c.children {
			child.HandleMessage(msg)
		}
	default:
		// For my children
		for _, child := range c.children {
			if child.HandleMessage(msg) {
				return true
			}
		}
	}

	return isStop
}

func (c *Component) reorderChildren() {
	sort.Sort(ByZorder(c.children))

	// Redraw me. Use message to refresh screen.
	SendMessage(BuildDrawMessage(c.Handler()))
}

func (c *Component) drawChildren() {
	for _, child := range c.children {
		child.HandleMessage(BuildDrawMessage(child.Handler()))
	}
}

// Name return the name of component.
func (c *Component) Name() string {
	return c.name
}

// Handler return handler value to send message to this component.
func (c *Component) Handler() uuid.UUID {
	return c.handler
}

// SetEnabled active or disable component.
func (c *Component) SetEnabled(status bool) {
	if c.OnEnabled != nil {
		c.enabled = c.OnEnabled(c, status)
	} else {
		c.enabled = status
	}
}

// GetEnabled return if component is enable.
func (c *Component) GetEnabled() bool {
	return c.enabled
}

// SetParent set parent of component. Nil -> no parent, maybe root component.
func (c *Component) SetParent(p TComponent) {
	c.parent = p
}

// GetParent return parent of component. Nil -> no parent, maybe root component.
func (c *Component) GetParent() TComponent {
	return c.parent
}

// AddChild add a child component.
func (c *Component) AddChild(child TComponent) {
	c.children = append(c.children, child)

	sort.Sort(ByZorder(c.children))
}

// RemoveChild remove the child.
func (c *Component) RemoveChild(child TComponent) {
	index := c.findChild(child.Handler())

	if index >= 0 {
		c.children[index].SetParent(nil)

		// Copy last element to index.
		c.children[index] = c.children[len(c.children)-1]
		// Truncate slice.
		c.children = c.children[:len(c.children)-1]
	}
}

// Find a child by uuid.
func (c *Component) findChild(uuid uuid.UUID) int {
	for i, child := range c.children {
		if child.Handler() == uuid {
			return i
		}
	}
	return -1
}

// Children return children list.
func (c *Component) Children() []TComponent {
	return c.children
}

// SetZorder order of message. Remember you must send WmZorderChange message
// to his parent.
func (c *Component) SetZorder(i int) {
	c.zorder = i
}

// GetZorder order of message.
func (c *Component) GetZorder() int {
	return c.zorder
}

// NewComponent create new component.
func NewComponent(name string) Component {
	return Component{
		name:    name,
		handler: uuid.New(),
	}
}
