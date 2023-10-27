package go_iterators

// ClosingIterator adds custom Closing logic on top of another iterator
type ClosingIterator[T any] struct {
	innerIterator Iterator[T]
	// extra function called AFTER "innerErr := Close()" returns
	close    func(innerErr error) error
	isClosed bool
}

func (c *ClosingIterator[T]) Next() (T, error) {
	return c.innerIterator.Next()
}

func (c *ClosingIterator[T]) Close() error {
	if c.isClosed {
		return ClosedIterator
	}
	err := c.innerIterator.Close()
	err = c.close(err)

	if err == nil {
		c.isClosed = true // if closing returned an error -> do not consider the iterator as closed.
	}

	return err
}

func NewClosingIterator[T any](innerIterator Iterator[T], close func(innerErr error) error) Iterator[T] {
	return &ClosingIterator[T]{
		innerIterator: innerIterator,
		close:         close,
	}
}
