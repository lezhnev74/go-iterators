package go_iterators

import "errors"

// SortedSelectingIterator returns sorted values from two other iterators
// if iterators are not sorted the behaviour is less predictable.
type SortedSelectingIterator[T any] struct {
	it1, it2             Iterator[T]
	v1, v2               T    // prefetched from internal iterators
	v1Fetched, v2Fetched bool // value exists
	cmp                  CmpFunc[T]
}

func (s *SortedSelectingIterator[T]) Close() error {
	err := s.it1.Close()
	if err != nil {
		s.it2.Close() // close anyway
		return err
	}
	return s.it2.Close()
}
func (si *SortedSelectingIterator[T]) fetch() error {
	var err error
	if !si.v1Fetched {
		si.v1, err = si.it1.Next()
		si.v1Fetched = err == nil
	}
	if err != nil && !errors.Is(err, EmptyIterator) {
		return err
	}
	if !si.v2Fetched {
		si.v2, err = si.it2.Next()
		si.v2Fetched = err == nil
	}
	if err != nil && !errors.Is(err, EmptyIterator) {
		return err
	}
	return nil
}

func (si *SortedSelectingIterator[T]) Next() (v T, err error) {
	err = si.fetch()
	if err != nil {
		return
	}

	if !si.v1Fetched && !si.v2Fetched {
		err = EmptyIterator
		return
	}

	// 1. only v1
	if si.v1Fetched && !si.v2Fetched {
		si.v1Fetched = false
		v = si.v1
		return
	}
	// 2. only v2
	if si.v2Fetched && !si.v1Fetched {
		si.v2Fetched = false
		v = si.v2
		return
	}
	// 3. both present
	r := si.cmp(si.v1, si.v2)
	if r == 0 {
		si.v1Fetched = false
		v = si.v1
		return
	} else if r < 0 {
		si.v1Fetched = false
		v = si.v1
		return
	} else {
		si.v2Fetched = false
		v = si.v2
		return
	}
}

func NewSortedSelectingIterator[T any](it1, it2 Iterator[T], cf CmpFunc[T]) Iterator[T] {
	return &SortedSelectingIterator[T]{
		it1: it1,
		it2: it2,
		cmp: cf,
	}
}
