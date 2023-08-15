package iterators

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

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

			s := NewDynamicSliceIterator(func() []any {
				if len(tt.sls) == 0 {
					return nil
				}
				sl := tt.sls[0]
				tt.sls = tt.sls[1:]
				return sl
			})
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
