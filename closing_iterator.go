package go_iterators

// ClosingIterator adds custom Closing logic on top of another iterator
type ClosingIterator[T any] struct {
	innerIterator Iterator[T]
	// extra function called AFTER "innerErr := Close()" returns
	close func(innerErr error) error
}

func (c *ClosingIterator[T]) Next() (T, error) {
	return c.innerIterator.Next()
}

func (c *ClosingIterator[T]) Close() error {
	err := c.innerIterator.Close()
	return c.close(err)
}

func NewClosingIterator[T any](innerIterator Iterator[T], close func(innerErr error) error) Iterator[T] {
	return &ClosingIterator[T]{
		innerIterator: innerIterator,
		close:         close,
	}
}
