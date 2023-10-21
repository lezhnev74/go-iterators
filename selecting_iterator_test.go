package go_iterators

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelectingIteratorClosesInnerIterators(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewSliceIterator([]string{})
	i3 := NewRemovingIterator[string](i1, i2, OrderedCmpFunc[string])
	require.NoError(t, i3.Close())
	require.ErrorIs(t, i3.Close(), ClosedIterator)
	require.ErrorIs(t, i1.Close(), ClosedIterator)
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestItClosesEmptyInnerIterator(t *testing.T) {
	i1 := NewSliceIterator([]string{"a"})
	i2 := NewSliceIterator([]string{})
	i3 := NewSortedSelectingIterator[string](i1, i2, OrderedCmpFunc[string])

	// Next() reads from it1 and gets EmptyIterator from it2 (thus closes it2)
	v, err := i3.Next()
	require.NoError(t, err)
	require.EqualValues(t, "a", v)

	err = i2.Close()
	require.ErrorIs(t, err, ClosedIterator)

	// Next() gets EmptyIterator from it1  (thus closes it1) and IteratorClosed from it2
	_, err = i3.Next()
	require.ErrorIs(t, err, EmptyIterator)

	err = i1.Close()
	require.ErrorIs(t, err, ClosedIterator)
	err = i2.Close()
	require.ErrorIs(t, err, ClosedIterator)
}
