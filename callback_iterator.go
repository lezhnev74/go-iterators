package iterators

type CallbackIterator[T any] func() (T, error)

func NewCallbackIterator[T any](cb func() (T, error)) Iterator[T] {
	var it CallbackIterator[T]
	it = cb
	return it
}

func (c CallbackIterator[T]) Close() error           { return nil }
func (c CallbackIterator[T]) Next() (v T, err error) { return c() }
