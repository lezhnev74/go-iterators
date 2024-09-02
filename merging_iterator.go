package go_iterators

import (
	"errors"
	"fmt"
	"slices"
)

// MergingIterator accepts N readers and outputs pending TermValues
// When all readers return EOF, so does this iterator.
type MergingIterator[T any] struct {
	// buf contains topmost value for each reader,
	// when reader is exhausted, it is removed from the buf
	buf []*IteratorCache[T]
	cmp CmpFunc[T]
	// merge accepts equal values and merges them into a new value
	merge    func(a, b T) T
	isClosed bool
}

func (mi *MergingIterator[T]) Next() (merged T, err error) {

	// 1. Fetch values from pending iterators
	err = mi.fetch()
	if err != nil {
		err = fmt.Errorf("merge iterator: %w", err)
		return
	}

	if len(mi.buf) == 0 {
		err = EmptyIterator
		return
	}

	// 2. Sort fetched values
	slices.SortFunc(mi.buf, func(a, b *IteratorCache[T]) int {
		return mi.cmp(a.v, b.v)
	})

	// 3. Merge equal values
	for i := range mi.buf {
		if i == 0 {
			merged = mi.buf[i].v // the first value goes out
			mi.buf[i].pending = true
			continue
		}

		if mi.cmp(merged, mi.buf[i].v) != 0 {
			break // a non-equal value will go out next time
		}
		merged = mi.merge(merged, mi.buf[i].v) // another equal value is pending
		mi.buf[i].pending = true
	}

	return
}

func (mi *MergingIterator[T]) Close() (err error) {
	if mi.isClosed {
		return ClosedIterator
	}
	for i, rc := range mi.buf {
		if rc == nil {
			continue
		}
		lastErr := rc.it.Close()
		if lastErr != nil && err == nil {
			err = lastErr // remember the first one
		}
		mi.buf[i] = nil // gc
	}
	mi.isClosed = true
	return
}

// fetch pulls data from each iterator which value was used in the merging,
// so the buffer contains the topmost value from each iterator.
func (mi *MergingIterator[T]) fetch() (err error) {
	var i int
	for j, rc := range mi.buf {
		if mi.buf[j] == nil {
			continue
		}

		mi.buf[j] = nil // gc

		if !rc.pending {
			mi.buf[i] = rc // keep in the buf
			i++
			continue
		}

		rc.v, err = rc.it.Next()
		if err == nil {
			rc.pending = false // just fetched
			mi.buf[i] = rc     // keep in the buf
			i++
			continue
		}
		if errors.Is(err, EmptyIterator) {
			err = rc.it.Close() // exhausted normally
			continue
		}
		return err // something bad happened in the underlying iterator
	}
	mi.buf = mi.buf[:i]
	return nil
}

func NewMergingIterator[T any](srcs []Iterator[T], cmpf CmpFunc[T], merge func(a, b T) T) *MergingIterator[T] {
	buf := make([]*IteratorCache[T], 0, len(srcs))
	for _, it := range srcs {
		buf = append(buf, &IteratorCache[T]{it: it, pending: true})
	}

	return &MergingIterator[T]{
		buf:   buf,
		cmp:   cmpf,
		merge: merge,
	}
}
