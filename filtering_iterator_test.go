package go_iterators

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFilteringIterator(t *testing.T) {
	inner := NewSliceIterator([]int{1, 2, 3, 4, 5})
	filter := func(i int) bool {
		return i%2 == 0
	}

	filteringIt := NewFilteringIterator(inner, filter)
	out := ToSlice(filteringIt)
	require.EqualValues(t, []int{2, 4}, out)

	require.NoError(t, filteringIt.Close())
	require.ErrorIs(t, filteringIt.Close(), ClosedIterator)
}
