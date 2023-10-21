package go_iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestItReturnsErrorFromInnerIterator_BatchingIterator(t *testing.T) {
	expectedError := fmt.Errorf("inner error")
	i1 := NewDynamicSliceIterator(
		func() ([]string, error) {
			return nil, expectedError
		},
		func() error { return nil },
	)
	i2 := NewBatchingIterator(i1, 10)
	_, err := i2.Next()
	require.ErrorIs(t, err, expectedError)
}

func TestItClosesInnerIterators_BatchingIterator(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewBatchingIterator(i1, 10)
	require.NoError(t, i1.Close())
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestItPanicsOnLowBatchSize(t *testing.T) {
	defer func() { _ = recover() }()
	NewBatchingIterator[string](nil, 0)
	t.Errorf("did not panic")
}

func TestBatchingIterator(t *testing.T) {
	type test struct {
		inner     []int
		batchSize int
		expected  [][]int
	}

	tests := []test{
		{[]int{}, 1, [][]int{}},
		{[]int{1}, 1, [][]int{{1}}},
		{[]int{1, 2}, 1, [][]int{{1}, {2}}},
		{[]int{1, 2, 3}, 2, [][]int{{1, 2}, {3}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewBatchingIterator(NewSliceIterator(tt.inner), tt.batchSize)
			result := ToSlice(s)

			if len(tt.expected) == 0 {
				require.Empty(t, result)
			} else {
				require.EqualValues(t, tt.expected, result)
			}
		})
	}
}
