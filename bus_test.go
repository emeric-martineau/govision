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
	"testing"

	"github.com/google/uuid"
)

// BroadcastHandler is use to send message to all components.

func TestBus_GetChannelMessage_and_SendMessage(t *testing.T) {
	m := BuildEmptyMessage()
	SendMessage(m)
}

func TestBus_BroadcastHandler(t *testing.T) {
	var broadcastHandler uuid.UUID = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	if BroadcastHandler() != broadcastHandler {
		t.Errorf("BroadcastHandler is change but I don't know !")
	}
}
