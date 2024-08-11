package go_iterators

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSequentialDynamicIterator(t *testing.T) {
	subIterators := []Iterator[int]{
		NewSliceIterator[int]([]int{1, 2}),
		NewSliceIterator[int]([]int{3, 4}),
	}
	next := func() (v Iterator[int], err error) {
		if len(subIterators) == 0 {
			return nil, EmptyIterator
		}
		v, subIterators = subIterators[0], subIterators[1:]
		return
	}

	it := NewSequentialDynamicIterator(next)
	values := []int{}
	for {
		v, err := it.Next()
		if errors.Is(err, EmptyIterator) {
			break
		}
		require.NoError(t, err)
		values = append(values, v)
	}
	require.Equal(t, []int{1, 2, 3, 4}, values)
}
