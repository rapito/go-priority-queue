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
	h []Interface
}

// New returns an initialized priority queue with the given elements.
// A call of the form New(x...) uses the underlying array of x to implement the queue
// and hence might change the elements of x.
// The complexity is O(n), where n = x.Len().
func New(x ...Interface) Queue {
	q := Queue{x}
	h := q.h
	for i := len(h) - 1; i >= 0; i-- {
		h[i].Index(i)
	}
	heapify(h)
	return q
}

// Push pushes the element x onto the queue.
// The complexity is O(log(n)) where n = q.Len().
func (q *Queue) Push(x Interface) {
	h := q.h
	n := len(h)
	q.h = append(h, x)
	up(q.h, n) // x.Index(n) is done by up.
}

// Pop removes a minimum element (according to Less) from the queue and returns it.
// The complexity is O(log(n)), where n = q.Len().
func (q *Queue) Pop() Interface {
	h := q.h
	n := len(h) - 1
	x := h[0]
	h[0], h[n] = h[n], h[0]
	down(h, 0, n) // h[0].Index(0) is done by down.
	h[n] = nil
	q.h = h[:n]
	return x
}

// Peek returns, but does not remove, a minimum element (according to Less) of the queue.
func (q *Queue) Peek() Interface {
	return q.h[0]
}

// Remove removes the element at index i from the queue and returns it.
// The complexity is O(log(n)), where n = q.Len().
func (q *Queue) Remove(i int) Interface {
	h := q.h
	n := len(h) - 1
	x := h[i]
	if n != i {
		h[i], h[n] = h[n], h[i]
		down(h, i, n) // h[i].Index(i) is done by down.
		up(h, i)
	}
	h[n] = nil
	q.h = h[:n]
	return x
}

// Len returns the number of elements in the queue.
func (q *Queue) Len() int {
	return len(q.h)
}

// Establishes the heap invariant in O(n) time.
func heapify(h []Interface) {
	n := len(h)
	for i := n/2 - 1; i >= 0; i-- {
		down(h, i, n)
	}
}

// Moves element at position j towards top of heap to restore invariant.
func up(h []Interface, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || h[i].Less(h[j]) {
			h[j].Index(j)
			break
		}
		h[i], h[j] = h[j], h[i]
		h[j].Index(j)
		j = i
	}
}

// Moves element at position i towards bottom of heap to restore invariant.
func down(h []Interface, i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n {
			h[i].Index(i)
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !h[j1].Less(h[j2]) {
			j = j2 // = 2*i + 2  // right child
		}
		if h[i].Less(h[j]) {
			h[i].Index(i)
			break
		}
		h[i], h[j] = h[j], h[i]
		h[i].Index(i)
		i = j
	}
}

// Returns the element at index i in the queue. Exported for testing.
func (q *Queue) get(i int) Interface {
	return q.h[i]
}
