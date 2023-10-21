package go_iterators

import (
	"cmp"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelectingIteratorClosesBothInnerIteratorsEvenIf1stErrors(t *testing.T) {
	badClose := fmt.Errorf("bad close")

	i1 := NewClosingIterator(
		NewSliceIterator([]string{}),
		func(innerError error) error { return badClose }, // always errors
	)
	i2 := NewSliceIterator([]string{})
	i3 := NewRemovingIterator[string](i1, i2, cmp.Compare[string])
	require.ErrorIs(t, i3.Close(), badClose)
	require.ErrorIs(t, i3.Close(), ClosedIterator)

	require.ErrorIs(t, i1.Close(), badClose)
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestSelectingIteratorClosesBothInnerIteratorsEvenIf2ndErrors(t *testing.T) {
	badClose := fmt.Errorf("bad close")

	i1 := NewSliceIterator([]string{})
	i2 := NewClosingIterator(
		NewSliceIterator([]string{}),
		func(innerError error) error { return badClose }, // always errors
	)
	i3 := NewRemovingIterator[string](i1, i2, cmp.Compare[string])
	require.ErrorIs(t, i3.Close(), badClose)
	require.ErrorIs(t, i3.Close(), ClosedIterator)

	require.ErrorIs(t, i1.Close(), ClosedIterator)
	require.ErrorIs(t, i2.Close(), badClose)
}

func TestSelectingIteratorClosesInnerIterators(t *testing.T) {
	i1 := NewSliceIterator([]string{})
	i2 := NewSliceIterator([]string{})
	i3 := NewRemovingIterator[string](i1, i2, cmp.Compare[string])
	require.NoError(t, i3.Close())
	require.ErrorIs(t, i3.Close(), ClosedIterator)
	require.ErrorIs(t, i1.Close(), ClosedIterator)
	require.ErrorIs(t, i2.Close(), ClosedIterator)
}

func TestItClosesEmptyInnerIterator(t *testing.T) {
	i1 := NewSliceIterator([]string{"a"})
	i2 := NewSliceIterator([]string{})
	i3 := NewSortedSelectingIterator[string](i1, i2, cmp.Compare[string])

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

func TestSelectingIteratorReturnsErrorFrom1stInnerIterator(t *testing.T) {
	badCall := fmt.Errorf("bad call")

	i1 := NewDynamicSliceIterator(
		func() ([]string, error) { return nil, badCall },
		func() error { return nil },
	)
	i2 := NewSliceIterator([]string{})
	i3 := NewRemovingIterator[string](i1, i2, cmp.Compare[string])

	_, err := i3.Next()
	require.ErrorIs(t, err, badCall)
}

func TestSelectingIteratorReturnsErrorFrom2ndInnerIterator(t *testing.T) {
	badCall := fmt.Errorf("bad call")

	i1 := NewSliceIterator([]string{})
	i2 := NewDynamicSliceIterator(
		func() ([]string, error) { return nil, badCall },
		func() error { return nil },
	)
	i3 := NewRemovingIterator[string](i1, i2, cmp.Compare[string])

	_, err := i3.Next()
	require.ErrorIs(t, err, badCall)
}
