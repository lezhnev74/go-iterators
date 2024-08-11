package go_iterators

import "errors"

// SequentialDynamicIterator will merge internal iterators into a single stream of values.
// nextIterator returns the next iterator to stream out, when next returns error, it is proxied out.
type SequentialDynamicIterator[T any] struct {
	Iterator[T]
	nextIterator func() (Iterator[T], error) // fetch the next iterator
	isClosed     bool
}

func (it *SequentialDynamicIterator[T]) Close() error {
	if it.isClosed {
		return ClosedIterator
	}
	it.isClosed = true

	if it.Iterator != nil {
		return it.Iterator.Close()
	}

	return nil
}

func (it *SequentialDynamicIterator[T]) Next() (v T, err error) {
	if it.Iterator == nil {
		it.Iterator, err = it.nextIterator()
		if err != nil {
			return
		}
	}

	v, err = it.Iterator.Next()
	if errors.Is(err, EmptyIterator) {
		err = it.Iterator.Close()
		if err != nil {
			return
		}
		it.Iterator = nil
		return it.Next()
	}
	return
}

func NewSequentialDynamicIterator[T any](nextIterator func() (Iterator[T], error)) *SequentialDynamicIterator[T] {
	return &SequentialDynamicIterator[T]{
		nextIterator: nextIterator,
	}
}
