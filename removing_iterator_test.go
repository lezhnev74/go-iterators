package lezhnev74

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRemIteratorClosesInnerIterators(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewSliceIterator([]string{})
	i3 := NewRemovingIterator[string](i1, i2, OrderedCmpFunc[string])
	require.NoError(t, i3.Close())
	require.ErrorIs(t, i3.Close(), ClosedIterator)
	require.ErrorIs(t, i1.Close(), ClosedIterator)
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestRemovingIterator(t *testing.T) {
	type test struct {
		sl1, sl2 []int
		f        CmpFunc[int]
		expected []int
	}
	fInt := func(a, b int) int {
		if a == b {
			return 0
		} else if a < b {
			return -1
		} else {
			return 1
		}
	}

	tests := []test{
		{sl1: []int{}, sl2: []int{}, f: fInt, expected: []int{}},
		{sl1: []int{1}, sl2: []int{}, f: fInt, expected: []int{1}},
		{sl1: []int{}, sl2: []int{1}, f: fInt, expected: []int{}},
		{sl1: []int{1, 2}, sl2: []int{1}, f: fInt, expected: []int{2}},
		{sl1: []int{1, 2, 3}, sl2: []int{1, 3}, f: fInt, expected: []int{2}},
		{sl1: []int{2}, sl2: []int{1, 3}, f: fInt, expected: []int{2}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewRemovingIterator[int](
				NewSliceIterator(tt.sl1),
				NewSliceIterator(tt.sl2),
				tt.f,
			)
			sl := make([]int, 0)
			for {
				v, err := s.Next()
				if err != nil {
					break
				}
				sl = append(sl, v)
			}
			require.EqualValues(t, tt.expected, sl)
		})
	}
}
