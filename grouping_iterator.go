package go_iterators

import (
	"errors"
)

// GroupingIterator batches values from the internal iterator based on criteria
type GroupingIterator[T any] struct {
	innerIterator Iterator[T]
	groupBy       func(T) any
	curBatch      []T
	curGroup      any
}

func NewGroupingIterator[T any](inner Iterator[T], groupBy func(T) any) Iterator[[]T] {
	return &GroupingIterator[T]{
		innerIterator: inner,
		groupBy:       groupBy,
	}
}

func (b *GroupingIterator[T]) Next() (v []T, err error) {

	var item T

	for {
		item, err = b.innerIterator.Next()

		if err != nil {
			break
		}

		itemGroup := b.groupBy(item)
		if itemGroup == b.curGroup || len(b.curBatch) == 0 {
			b.curBatch = append(b.curBatch, item)
			b.curGroup = itemGroup
			continue
		}

		v = b.curBatch
		b.curBatch = []T{item}
		b.curGroup = itemGroup
		return
	}

	if errors.Is(err, EmptyIterator) && len(v) > 0 {
		err = nil
	}

	if len(b.curBatch) > 0 {
		err = nil
		v = b.curBatch
		b.curBatch = nil
		return
	}

	return
}

func (b *GroupingIterator[T]) Close() error {
	return b.innerIterator.Close()
}
