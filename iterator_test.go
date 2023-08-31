package go_iterators

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToSlice(t *testing.T) {
	sl := []int{1, 2, 3}
	s := NewSliceIterator(sl)
	sl1 := ToSlice[int](s)
	require.EqualValues(t, sl, sl1)
}
