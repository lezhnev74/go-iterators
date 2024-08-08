package go_iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupingIterator(t *testing.T) {
	type test struct {
		inner    []int
		groupBy  func(int) any
		expected [][]int
	}

	odd := func(i int) any {
		if i%2 == 0 {
			return "T"
		}
		return "F"
	}
	less3 := func(i int) any {
		if i < 3 {
			return "T"
		}
		return "F"
	}

	tests := []test{
		{[]int{}, odd, [][]int{}},
		{[]int{1}, odd, [][]int{{1}}},
		{[]int{1, 2}, odd, [][]int{{1}, {2}}},
		{[]int{1, 2, 3}, odd, [][]int{{1}, {2}, {3}}},
		{[]int{1, 2, 3}, less3, [][]int{{1, 2}, {3}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewGroupingIterator(NewSliceIterator(tt.inner), tt.groupBy)
			result := ToSlice(s)

			if len(tt.expected) == 0 {
				require.Empty(t, result)
			} else {
				require.EqualValues(t, tt.expected, result)
			}
		})
	}
}

func TestItClosesInnerIterators_GroupingIterator(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewGroupingIterator(i1, func(s string) any { return s })
	require.NoError(t, i1.Close())
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestItReturnsErrorFromInnerIterator_GroupingIterator(t *testing.T) {
	expectedError := fmt.Errorf("inner error")
	i1 := NewDynamicSliceIterator(
		func() ([]string, error) {
			return nil, expectedError
		},
		func() error { return nil },
	)
	i2 := NewGroupingIterator(i1, func(s string) any { return s })
	_, err := i2.Next()
	require.ErrorIs(t, err, expectedError)
}
