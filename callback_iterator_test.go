package iterators

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCallbackIterator(t *testing.T) {
	src := []string{"a", "b", "c"}
	cb := func() (string, error) {
		if len(src) > 0 {
			s := src[0]
			src = src[1:]
			return s, nil
		}
		return "", EmptyIterator
	}
	it := NewCallbackIterator(cb)
	require.EqualValues(t, []string{"a", "b", "c"}, ToSlice(it))
}
