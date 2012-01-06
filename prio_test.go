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
	. "prio"
	"testing"
)

type myInt int

func (x myInt) Less(y Interface) bool { return x < y.(myInt) }
func (x myInt) Index(i int)           {}

type myType struct {
	value int
	index int // index in heap
}

func (x *myType) Less(y Interface) bool { return x.value < y.(*myType).value }
func (x *myType) Index(i int)           { x.index = i }

// Verify the ordering of all elements in the heap.
func verify(t *testing.T, q Queue) {
	n := q.Len()
	for i := 1; i < n; i++ {
		p := (i - 1) / 2 // parent
		qi := q.Get(i)
		qp := q.Get(p)
		if qi.Less(qp) {
			t.Errorf("heap invariant invalidated [%d] = %v < [%d] = %v", i, qi, p, qp)
		}
	}
}

// Verify that all elements in a queue of myType have been given the correct index.
func verifyIndex(t *testing.T, q Queue) {
	for i := 0; i < q.Len(); i++ {
		index := q.Get(i).(*myType).index
		if index != i {
			t.Errorf("wrong index [%d] = %d", i, index)
		}
	}
}

func TestNew0(t *testing.T) {
	a := make([]Interface, 10)
	for i := 0; i < 10; i++ {
		a[i] = myInt(0)
	}
	q := New(a...)
	verify(t, q)

	for i := 1; q.Len() > 0; i++ {
		x := q.Pop().(myInt)
		verify(t, q)
		if x != 0 {
			t.Errorf("%d.th pop got %d; want %d", i, x, 0)
		}
	}
}

func TestNew1(t *testing.T) {
	a := make([]Interface, 10)
	for i := 0; i < 10; i++ {
		a[i] = myInt(i + 1)
	}
	q := New(a...)
	verify(t, q)

	for i := 1; q.Len() > 0; i++ {
		x := q.Pop().(myInt)
		verify(t, q)
		if int(x) != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func Test(t *testing.T) {
	q := New()
	for i := 20; i > 10; i-- {
		q.Push(myInt(i))
		verify(t, q)
	}

	for i := 10; i > 0; i-- {
		q.Push(myInt(i))
		verify(t, q)
	}

	for i := 1; q.Len() > 0; i++ {
		x := q.Peek().(myInt)
		y := q.Pop().(myInt)
		verify(t, q)
		if i < 20 {
			q.Push(myInt(20 + i))
			verify(t, q)
		}
		if int(x) != i {
			t.Errorf("%d.th peek got %d; want %d", i, x, i)
		}
		if int(y) != i {
			t.Errorf("%d.th pop got %d; want %d", i, y, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	a := make([]*myType, 10)
	q := Queue{}
	for i := 0; i < len(a); i++ {
		a[i] = &myType{i, 0}
		q.Push(a[i])
		verify(t, q)
		verifyIndex(t, q)
	}

	for i := 0; i < len(a); i++ {
		x := q.Remove(0)
		if x != a[i] {
			t.Errorf("Remove(0) got %v; want %v", x, a[i])
		}
		verify(t, q)
		verifyIndex(t, q)
	}
}

func TestRemove1(t *testing.T) {
	a := make([]*myType, 10)
	q := Queue{}
	for i := 0; i < len(a); i++ {
		a[i] = &myType{i, 0}
		q.Push(a[i])
		verify(t, q)
		verifyIndex(t, q)
	}

	for i := len(a) - 1; i >= 0; i-- {
		index := a[i].index
		x := q.Remove(index)
		if x != a[i] {
			t.Errorf("Remove(%d) got %v; want %v", index, x, a[i])
		}
		verify(t, q)
		verifyIndex(t, q)
	}
}

func TestRemove2(t *testing.T) {
	a := make([]Interface, 10)
	for i := len(a) - 1; i >= 0; i-- {
		a[i] = &myType{i, 0}
	}
	q := New(a...)
	verify(t, q)
	verifyIndex(t, q)

	for i := len(a) - 1; i >= 0; i-- {
		x := a[i]
		index := x.(*myType).index
		y := q.Remove(index)
		if x != y {
			t.Errorf("Remove(%d) got %v; want %v", index, y, x)
		}
		verify(t, q)
		verifyIndex(t, q)
	}
}
