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

// ByZorder implements sort.Interface based on the zorder field.
type ByZorder []TComponent

func (a ByZorder) Len() int {
	return len(a)
}

func (a ByZorder) Less(i, j int) bool {
	return a[i].GetZorder() < a[j].GetZorder()
}

func (a ByZorder) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
