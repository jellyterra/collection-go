// Copyright 2024 Jelly Terra
// Use of this source code is governed by the MIT license.

package collection

import "sync"

func ValuesOfRawMap[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i++
	}
	return values
}

func MergeRawMap[K comparable, V any](set map[K]V, subset map[K]V) {
	for k, v := range subset {
		set[k] = v
	}
}

// Map provides wrapped Go map for various operations.
type Map[K comparable, V any] struct {
	Raw map[K]V
}

func (m *Map[K, V]) Len() int { return len(m.Raw) }

func (m *Map[K, V]) Foreach(f func(K, V) error) error {
	for k, v := range m.Raw {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (m *Map[K, V]) Delete(k K) { delete(m.Raw, k) }

func (m *Map[K, V]) Set(k K, v V) {
	if m.Raw == nil {
		m.Raw = make(map[K]V)
	}

	m.Raw[k] = v
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	v, ok := m.Raw[k]
	return v, ok
}

// SyncMap provides Map with RW-mutex protected.
type SyncMap[K comparable, V any] struct {
	It      Map[K, V]
	RWMutex sync.RWMutex
}

func (s *SyncMap[K, V]) Do(f func(it *Map[K, V])) {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	f(&s.It)
}

func (s *SyncMap[K, V]) Get(k K) (V, bool) {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.It.Get(k)
}

func (s *SyncMap[K, V]) Foreach(f func(K, V) error) error {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.It.Foreach(f)
}

func (s *SyncMap[K, V]) DoMut(f func(it *Map[K, V])) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	f(&s.It)
}

func (s *SyncMap[K, V]) Delete(k K) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	s.It.Delete(k)
}

func (s *SyncMap[K, V]) Set(k K, v V) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	s.It.Set(k, v)
}
