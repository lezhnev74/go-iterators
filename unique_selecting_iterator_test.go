package go_iterators

import (
	"cmp"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestItStopsIfInnerIteratorErrors_UniqSelect(t *testing.T) {
	badCall := fmt.Errorf("bad call")

	i1 := NewSliceIterator([]string{})
	i2 := NewDynamicSliceIterator(
		func() ([]string, error) { return nil, badCall },
		func() error { return nil },
	)
	i3 := NewUniqueSelectingIterator[string](i1, i2, cmp.Compare[string])

	_, err := i3.Next()
	require.ErrorIs(t, err, badCall)
}

func TestItClosesInnerIterators(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewSliceIterator([]string{})
	i3 := NewUniqueSelectingIterator[string](i1, i2, cmp.Compare[string])
	require.NoError(t, i3.Close())
	require.ErrorIs(t, i3.Close(), ClosedIterator)
	require.ErrorIs(t, i1.Close(), ClosedIterator)
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestUniqueSelectingIterator(t *testing.T) {
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
		{sl1: []int{}, sl2: []int{1}, f: fInt, expected: []int{1}},
		{sl1: []int{1}, sl2: []int{}, f: fInt, expected: []int{1}},
		{sl1: []int{1, 2}, sl2: []int{}, f: fInt, expected: []int{1, 2}},
		{sl1: []int{1}, sl2: []int{1}, f: fInt, expected: []int{1}},
		{sl1: []int{1, 2}, sl2: []int{1}, f: fInt, expected: []int{1, 2}},
		{sl1: []int{1, 2}, sl2: []int{1, 2}, f: fInt, expected: []int{1, 2}},
		{sl1: []int{1, 2}, sl2: []int{2, 3}, f: fInt, expected: []int{1, 2, 3}},
		{sl1: []int{1, 4}, sl2: []int{2, 3}, f: fInt, expected: []int{1, 2, 3, 4}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewUniqueSelectingIterator[int](
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
