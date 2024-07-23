package go_iterators

import (
	"errors"
	"io"
)

// EmptyIterator shows that no further value will follow from the iterator
var EmptyIterator = errors.New("iterator is empty")
var ClosedIterator = errors.New("iterator is closed")

// Iterator is used for working with sequences of possibly unknown size
// Interface adds a performance penalty for indirection.
type Iterator[T any] interface {
	// Next returns EmptyIterator when no value available at the source
	// error == nil means the returned value is good
	Next() (T, error)
	// Closer the client may decide to stop the iteration before EmptyIterator received
	// After the first call it must return ClosedIterator.
	io.Closer
}

// CmpFunc returns -1,0,1 respectively if a<b,a=b,a>b
type CmpFunc[T any] func(a, b T) int

func ToSlice[T any](it Iterator[T]) (dump []T) {
	for {
		v, err := it.Next()
		if err != nil {
			break
		}
		dump = append(dump, v)
	}
	return
}

// IteratorCache keeps the topmost value from the iterator.
type IteratorCache[T any] struct {
	// internal wrapped iterator
	it Iterator[T]
	// topmost fetched value
	v T
	// when the value is merged, this flag is set to signify the need of a fetch
	pending bool
}
