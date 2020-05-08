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
	"github.com/google/uuid"
)

// Is use to send message to all components.
var broadcastHandler uuid.UUID = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

// Is use to send message to application only.
var applicationHandler uuid.UUID = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

// Bus bus message.
type Bus struct {
	busChannel chan Message
}

// Send a event in bus.
func (b Bus) Send(e Message) {
	b.busChannel <- e
}

// Channel a event in bus.
func (b Bus) Channel() *chan Message {
	return &b.busChannel
}

// BroadcastHandler return value for broadcast all component.
func BroadcastHandler() uuid.UUID {
	return broadcastHandler
}

// ApplicationHandler is use to send message to application only.
func ApplicationHandler() uuid.UUID {
	return applicationHandler
}

// NewBus create a new bus.
func NewBus() Bus {
	return Bus{
		busChannel: make(chan Message, 10),
	}
}
