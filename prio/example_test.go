// Copyright 2012 Stefan Nilsson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prio_test

import (
	prio "code.google.com/p/go-priority-queue/prio"
	"fmt"
)

type myType struct {
	value int
	index int // index in heap
}

func (x *myType) Less(y prio.Interface) bool { return x.value < y.(*myType).value }
func (x *myType) Index(i int)                { x.index = i }

func ExampleQueue() {
	var q prio.Queue

	// Create and push five elements with values 1 to 5 onto the queue.
	a := make([]*myType, 5)
	for i := range a {
		a[i] = &myType{value: i + 1}
		q.Push(a[i])
	}

	q.Remove(a[1].index) // Use index to find and remove second element.
	for q.Len() > 0 {
		fmt.Print(q.Pop().(*myType).value)
	}
	// Output: 1345
}
