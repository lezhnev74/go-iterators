package go_iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClosingIterator(t *testing.T) {
	var closed int
	innerErr := fmt.Errorf("inner error")

	data := []any{1}
	inner := NewCallbackIterator(func() (any, error) {
		if len(data) != 0 {
			a := data[0]
			data = data[1:]
			return a, nil
		}
		return nil, EmptyIterator
	}, func() error {
		closed++
		return innerErr
	})

	closing := NewClosingIterator(inner, func(innerErr error) error {
		closed++
		return innerErr
	})

	// Make sure inner iterator is proxied
	s := ToSlice(closing)
	require.EqualValues(t, []any{1}, s)

	// Make sure inner close is called too and the inner error is propagated
	err := closing.Close()
	require.Equal(t, err, innerErr)
	require.Equal(t, 2, closed)
}
