package go_iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMappingIterator(t *testing.T) {
	inner := NewSliceIterator([]int{1, 2, 3, 4, 5})
	mp := func(i int) string {
		return fmt.Sprintf("%d", i*10)
	}

	mIt := NewMappingIterator(inner, mp)
	out := ToSlice(mIt)
	require.EqualValues(t, []string{"10", "20", "30", "40", "50"}, out)

	require.NoError(t, mIt.Close())
	require.ErrorIs(t, mIt.Close(), ClosedIterator)
}
