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

// Package prio provides a priority queue.
//
// The queue can hold elements that implement the two methods of prio.Interface.
// The simplest use case looks like this:
//
//	type myInt int
//
//	func (x myInt) Less(y Interface) bool { return x < y.(myInt) }
//	func (x myInt) Index(i int)           {}
//
// To use the Remove method you need to keep track of the index of elements in the heap, 
// e.g. like this:
//
//	type myType struct {
//		value int
//		index int // index in heap
// 	}
//
//	func (x *myType) Less(y Interface) bool { return x.value < y.(*myType).value }
//	func (x *myType) Index(i int)           { x.index = i }
//
package prio

import "container/heap"

// A type that implements prio.Interface can be inserted into a priority queue.
type Interface interface {
	// Less returns whether this element should sort before element x.
	Less(x Interface) bool
	// Index is called by the priority queue when this element is moved to index i.
	Index(i int)
}

// Queue represents a priority queue.
// The zero value for Queue is an empty queue ready to use.
type Queue struct {
	h myHeap // a slice implementing heap.Interface
}

// New returns an initialized priority queue with the given elements.
// A call of the form New(x...) uses the underlying array of x to implement the queue
// and hence might change the elements of x.
// The complexity is O(n), where n = x.Len().
func New(x ...Interface) Queue {
	q := Queue{h: x}
	for i, v := range q.h {
		v.Index(i)
	}
	heap.Init(&q.h)
	return q
}

// Push pushes the element x onto the queue.
// The complexity is O(log(n)) where n = q.Len().
func (q *Queue) Push(x Interface) {
	heap.Push(&q.h, x)
}

// Pop removes a minimum element (according to Less) from the queue and returns it.
// The complexity is O(log(n)), where n = q.Len().
func (q *Queue) Pop() Interface {
	return heap.Pop(&q.h).(Interface)
}

// Peek returns, but does not remove, a minimum element (according to Less) of the queue.
func (q *Queue) Peek() Interface {
	return q.h[0]
}

// Remove removes the element at index i from the queue and returns it.
// The complexity is O(log(n)), where n = q.Len().
func (q *Queue) Remove(i int) Interface {
	return heap.Remove(&q.h, i).(Interface)
}

// Len returns the number of elements in the queue.
func (q *Queue) Len() int {
	return len(q.h)
}

type myHeap []Interface // implements heap.Interface

func (h myHeap) Len() int {
	return len(h)
}

func (h myHeap) Less(i, j int) bool {
	return h[i].Less(h[j])
}

func (h myHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index(i)
	h[j].Index(j)
}

func (h *myHeap) Push(x interface{}) {
	*h = append(*h, x.(Interface))
	x.(Interface).Index(len(*h) - 1)
}

func (h *myHeap) Pop() interface{} {
	a := *h
	i := len(a) - 1
	x := a[i]
	a[i] = nil
	*h = a[:i]
	return x
}

// Returns the element at index i in the queue. Exported for testing.
func (q *Queue) get(i int) Interface {
	return q.h[i]
}
