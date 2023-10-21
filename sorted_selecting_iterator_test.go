package go_iterators

import (
	"cmp"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestItClosesInnerIterators_SortedSelector(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewSliceIterator([]string{})
	i3 := NewSortedSelectingIterator[string](i1, i2, OrderedCmpFunc[string])
	require.NoError(t, i3.Close())
	require.ErrorIs(t, i3.Close(), ClosedIterator)
	require.ErrorIs(t, i1.Close(), ClosedIterator)
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestSortedSelectingIterator(t *testing.T) {
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
		{sl1: []int{1}, sl2: []int{1}, f: fInt, expected: []int{1, 1}},
		{sl1: []int{1, 2}, sl2: []int{1}, f: fInt, expected: []int{1, 1, 2}},
		{sl1: []int{1, 2}, sl2: []int{1, 2}, f: fInt, expected: []int{1, 1, 2, 2}},
		{sl1: []int{1, 2}, sl2: []int{2, 3}, f: fInt, expected: []int{1, 2, 2, 3}},
		{sl1: []int{1, 4}, sl2: []int{2, 3}, f: fInt, expected: []int{1, 2, 3, 4}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewSortedSelectingIterator[int](
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

func TestNestedSortedSelectingIterator(t *testing.T) {
	type test struct {
		sl1, sl2 Iterator[int]
		expected []int
	}

	tests := []test{
		{
			sl1:      NewSliceIterator([]int{1}),
			sl2:      NewSliceIterator([]int{2}),
			expected: []int{1, 2},
		},
		{
			sl1: NewSortedSelectingIterator(
				NewSliceIterator([]int{1}),
				NewSliceIterator([]int{3}),
				cmp.Compare[int],
			),
			sl2:      NewSliceIterator([]int{2}),
			expected: []int{1, 2, 3},
		},
		{
			sl1:      NewSliceIterator([]int{9, 10}),
			sl2:      NewSliceIterator([]int{9, 10, 11}),
			expected: []int{9, 9, 10, 10, 11},
		},
		{
			sl1:      NewSliceIterator([]int{1, 2, 3}),
			sl2:      NewSliceIterator([]int{1, 3}),
			expected: []int{1, 1, 2, 3, 3},
		},
		{
			sl1: NewSortedSelectingIterator(
				NewSliceIterator([]int{1}),
				NewSliceIterator([]int{3}),
				cmp.Compare[int],
			),
			sl2: NewSortedSelectingIterator(
				NewSliceIterator([]int{2}),
				NewSliceIterator([]int{1}),
				cmp.Compare[int],
			),
			expected: []int{1, 1, 2, 3},
		},
		{
			sl1: NewSortedSelectingIterator(
				NewSliceIterator([]int{1, 5, 9}),
				NewSliceIterator([]int{3}),
				cmp.Compare[int],
			),
			sl2: NewSortedSelectingIterator(
				NewSliceIterator([]int{2}),
				NewSliceIterator([]int{}),
				cmp.Compare[int],
			),
			expected: []int{1, 2, 3, 5, 9},
		},
		{
			sl1: NewSortedSelectingIterator(
				NewSortedSelectingIterator(
					NewSliceIterator([]int{5}),
					NewSliceIterator([]int{0}),
					cmp.Compare[int],
				),
				NewSliceIterator([]int{3}),
				cmp.Compare[int],
			),
			sl2: NewSortedSelectingIterator(
				NewSliceIterator([]int{2}),
				NewSliceIterator([]int{1}),
				cmp.Compare[int],
			),
			expected: []int{0, 1, 2, 3, 5},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewSortedSelectingIterator[int](tt.sl1, tt.sl2, cmp.Compare[int])
			require.EqualValues(t, tt.expected, ToSlice(s))
		})
	}
}
