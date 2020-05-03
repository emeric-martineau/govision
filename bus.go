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

var busChannel chan Message

// Is use to send message to all components.
var broadcastHandler uuid.UUID = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

// Is use to send message to application only.
var nullHandler uuid.UUID = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

// Inits a new bus.
func init() {
	busChannel = make(chan Message, 10)
}

// SendMessage a event in bus.
func SendMessage(e Message) {
	busChannel <- e
}

// BroadcastHandler return value for broadcast all component.
func BroadcastHandler() uuid.UUID {
	return broadcastHandler
}

// NullHandler is use to send message to application only.
func NullHandler() uuid.UUID {
	return nullHandler
}
