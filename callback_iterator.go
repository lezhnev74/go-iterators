package go_iterators

type CallbackIterator[T any] struct {
	cb       func() (T, error)
	close    func() error
	isClosed bool
}

func NewCallbackIterator[T any](
	cb func() (T, error),
	close func() error,
) Iterator[T] {
	return &CallbackIterator[T]{
		cb:    cb,
		close: close,
	}
}

func (c *CallbackIterator[T]) Close() error {
	if c.isClosed {
		return ClosedIterator
	}
	c.isClosed = true
	return c.close()
}
func (c *CallbackIterator[T]) Next() (v T, err error) { return c.cb() }
