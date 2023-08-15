package lezhnev74

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClose(t *testing.T) {
	s := NewSliceIterator([]string{})
	require.NoError(t, s.Close())
	require.ErrorIs(t, s.Close(), ClosedIterator)
}

func TestSliceIterator(t *testing.T) {
	type test struct {
		sl []any
	}
	tests := []test{
		{sl: []any{}},
		{sl: []any{"a", "b", "c"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			s := NewSliceIterator(tt.sl)
			sl2 := make([]any, 0)
			for {
				v, err := s.Next()
				if err != nil {
					break
				}
				sl2 = append(sl2, v)
			}
			require.EqualValues(t, tt.sl, sl2)
		})
	}
}
