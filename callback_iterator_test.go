package go_iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCallbackError(t *testing.T) {
	e := fmt.Errorf("callback failed")
	s := NewCallbackIterator(
		func() (any, error) { return nil, e },
		func() error { return nil },
	)

	_, err := s.Next()
	require.ErrorIs(t, err, e)
}

func TestCloseCallbackIterator(t *testing.T) {
	var closed bool
	s := NewCallbackIterator(
		func() (string, error) { return "a", nil },
		func() error {
			closed = true
			return nil
		},
	)
	require.NoError(t, s.Close())
	require.True(t, closed)
	require.ErrorIs(t, s.Close(), ClosedIterator)
}

func TestCallbackIterator(t *testing.T) {
	src := []string{"a", "b", "c"}
	cb := func() (string, error) {
		if len(src) > 0 {
			s := src[0]
			src = src[1:]
			return s, nil
		}
		return "", EmptyIterator
	}
	it := NewCallbackIterator(cb, func() error { return nil })
	require.EqualValues(t, []string{"a", "b", "c"}, ToSlice(it))
}
