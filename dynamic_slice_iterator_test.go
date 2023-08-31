package go_iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestErrorFetch(t *testing.T) {
	e := fmt.Errorf("fetch failed")
	s := NewDynamicSliceIterator(
		func() ([]string, error) { return nil, e },
		func() error { return nil },
	)

	_, err := s.Next()
	require.ErrorIs(t, err, e)
}

func TestCloseDynamicSliceIterator(t *testing.T) {
	var closed bool
	s := NewDynamicSliceIterator(
		func() ([]string, error) { return []string{"a"}, nil },
		func() error {
			closed = true
			return nil
		},
	)
	require.NoError(t, s.Close())
	require.True(t, closed)
	require.ErrorIs(t, s.Close(), ClosedIterator)
}

func TestDynamicSliceIterator(t *testing.T) {
	type test struct {
		sls [][]any
	}
	tests := []test{
		{sls: [][]any{}},
		{sls: [][]any{{"a", "b"}, {"c"}}},
		{sls: [][]any{{"a"}, {"c"}}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			expectedMergedSlice := make([]any, 0)
			for _, sl := range tt.sls {
				expectedMergedSlice = append(expectedMergedSlice, sl...)
			}

			s := NewDynamicSliceIterator(func() ([]any, error) {
				if len(tt.sls) == 0 {
					return nil, EmptyIterator
				}
				sl := tt.sls[0]
				tt.sls = tt.sls[1:]
				return sl, nil
			}, func() error { return nil })
			sl2 := make([]any, 0)
			for {
				v, err := s.Next()
				if err != nil {
					break
				}
				sl2 = append(sl2, v)
			}

			require.EqualValues(t, expectedMergedSlice, sl2)
		})
	}
}
