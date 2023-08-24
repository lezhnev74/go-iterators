package lezhnev74

import "errors"

// RemovingIterator returns all from it1 that are not present in it2
type RemovingIterator[T any] struct {
	it1, it2             Iterator[T]
	v1, v2               T    // prefetched from internal iterators
	v1Fetched, v2Fetched bool // is value prefetched
	cmp                  CmpFunc[T]
}

func (s *RemovingIterator[T]) Close() error {
	err := s.it1.Close()
	if err != nil {
		s.it2.Close() // close anyway
		return err
	}
	return s.it2.Close()
}
func (si *RemovingIterator[T]) fetch() error {
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

	// check if both values present and collapse, fetch more if so
	if si.v1Fetched && si.v2Fetched {
		r := si.cmp(si.v1, si.v2)
		if r == 0 {
			si.v1Fetched, si.v2Fetched = false, false
			err = si.fetch()
		}
	}
	if err != nil && !errors.Is(err, EmptyIterator) {
		return err
	}

	return nil
}

func (si *RemovingIterator[T]) Next() (v T, err error) {
	err = si.fetch()
	if err != nil {
		return
	}

	if si.v1Fetched {
		si.v1Fetched = false
		v = si.v1
		return
	}

	err = EmptyIterator
	return
}

func NewRemovingIterator[T any](itMain, itRemove Iterator[T], cf CmpFunc[T]) Iterator[T] {
	return &RemovingIterator[T]{
		it1: itMain,
		it2: itRemove,
		cmp: cf,
	}
}
