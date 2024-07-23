// Copyright 2024 Jelly Terra
// Use of this source code is governed by the MIT license.

package collection

import "sync"

// Vector provides wrapped Go slice for various operations.
type Vector[E any] struct {
	Raw []E
}

func (v *Vector[E]) Len() int { return len(v.Raw) }

func (v *Vector[E]) Cap() int { return cap(v.Raw) }

func (v *Vector[E]) Append(e ...E) Vector[E] { return Vector[E]{Raw: append(v.Raw, e...)} }

func (v *Vector[E]) Push(e ...E) { v.Raw = append(v.Raw, e...) }

func (v *Vector[E]) PopN(n int) { v.Raw = v.Raw[:v.Len()-n] }

func (v *Vector[E]) Pop() { v.Raw = v.Raw[:v.Len()-1] }

type ComparableVector[E comparable] struct {
	Vector[E]
}

func (v *ComparableVector[E]) Contains(e E) bool {
	for _, each := range v.Raw {
		if each == e {
			return true
		}
	}

	return false
}

// SyncVector provides Vector with RW-mutex protected.
type SyncVector[E any] struct {
	It Vector[E]

	RWMutex sync.RWMutex
}

// Do locks R-mutex during execution.
func (v *SyncVector[E]) Do(f func(it *Vector[E])) {
	v.RWMutex.RLock()
	defer v.RWMutex.RUnlock()
	f(&v.It)
}

func (v *SyncVector[E]) Append(e ...E) Vector[E] {
	v.RWMutex.RLock()
	defer v.RWMutex.RUnlock()
	return v.It.Append(e...)
}

// DoMut locks W-mutex during execution.
func (v *SyncVector[E]) DoMut(f func(it *Vector[E])) {
	v.RWMutex.Lock()
	defer v.RWMutex.Unlock()
	f(&v.It)
}

func (v *SyncVector[E]) Push(e ...E) {
	v.RWMutex.Lock()
	defer v.RWMutex.Unlock()
	v.It.Push(e...)
}

type SyncComparableVector[E comparable] struct {
	ComparableVector[E]

	lock sync.RWMutex
}

func (a *SyncComparableVector[E]) Contains(e E) bool {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.ComparableVector.Contains(e)
}
