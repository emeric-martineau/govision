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
	"testing"
)

func compare(c, r TComponent, t *testing.T) {
	if c != r {
		t.Errorf("Must be '%s' found '%s'", c.Name(), r.Name())
	}
}

// Test if function works :)
func TestSortComponent_ByZorder(t *testing.T) {
	c1 := NewComponent("1")
	c1.SetZorder(2)
	c2 := NewComponent("2")
	c2.SetZorder(1)
	c3 := NewComponent("3")
	c3.SetZorder(3)

	var children []TComponent

	children = append(children, &c1)
	children = append(children, &c2)
	children = append(children, &c3)

	sort.Sort(ByZorder(children))

	compare(children[0], &c2, t)
	compare(children[1], &c1, t)
	compare(children[2], &c3, t)
}
