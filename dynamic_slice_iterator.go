package iterators

// DynamicSliceIterator implements Iterator over a dynamic slice
// whenever it needs data it calls fetch() for a new slice to iterate
type DynamicSliceIterator[T any] struct {
	values []T
	fetch  func() []T // nil or empty slice leads to stopping iteration
}

func (s *DynamicSliceIterator[T]) Close() error { return nil }
func (s *DynamicSliceIterator[T]) Next() (v T, err error) {
	if len(s.values) == 0 {
		s.values = s.fetch()
	}

	if len(s.values) == 0 {
		err = EmptyIterator
		return
	}

	v = s.values[0]
	s.values = s.values[1:]
	return
}

func NewDynamicSliceIterator[T any](fetch func() []T) Iterator[T] {
	return &DynamicSliceIterator[T]{
		fetch: fetch,
	}
}
