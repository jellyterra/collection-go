// Copyright 2024 Jelly Terra
// Use of this source code is governed by the MIT license.

package collection

type Synchronized[T any] interface {
	Do(f func(it *T))
	DoMut(f func(it *T))
}
