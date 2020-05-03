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

// Ok, this test is only to don't have red line in IDE :)
func TestMessage_Messages(t *testing.T) {
	uuid := uuid.New()
	s := mkTestScreen(t, "")
	defer s.Fini()

	BuildEmptyMessage()
	BuildDrawMessage(uuid)
	BuildKeyMessage(nil)
	BuildChangeBounds(uuid, Rect{})
	BuildZorderMessage(uuid)
	BuildScreenResizeMessage(s)
}
