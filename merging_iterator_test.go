package go_iterators

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

type TermValues struct {
	term     []byte
	postings []uint64
}

func CompareTermValues(a, b TermValues) int {
	return bytes.Compare(a.term, b.term)
}

func MergeTermValues(a, b TermValues) (merged TermValues) {
	merged.term = a.term
	merged.postings = append(a.postings, b.postings...)

	return
}

func TestMergingIterator(t *testing.T) {
	input := [][]TermValues{
		// source 1
		{
			{[]byte("term1"), []uint64{10, 500, 30}},
			{[]byte("term2"), []uint64{1}},
		},
		// source 2
		{
			{[]byte("term2"), []uint64{99, 1}},
			{[]byte("term3"), []uint64{33}},
		},
		// source 3
		{
			{[]byte("term1a"), []uint64{0}},
			{[]byte("term2"), []uint64{5513}},
			{[]byte("term4"), []uint64{987, 11}},
		},
	}
	expected := []TermValues{
		{[]byte("term1"), []uint64{10, 30, 500}},
		{[]byte("term1a"), []uint64{0}},
		{[]byte("term2"), []uint64{1, 99, 5513}},
		{[]byte("term3"), []uint64{33}},
		{[]byte("term4"), []uint64{11, 987}},
	}

	readers := make([]Iterator[TermValues], 0, len(input))
	for _, tvs := range input {
		r := NewSliceIterator(tvs)
		readers = append(readers, r)
	}

	mi := NewMergingIterator[TermValues](readers, CompareTermValues, MergeTermValues)
	actual := make([]TermValues, 0)
	for {
		tv, err := mi.Next()
		if errors.Is(err, EmptyIterator) {
			break
		}
		require.NoError(t, err)
		actual = append(actual, tv)
	}
	require.NoError(t, mi.Close()) // todo test all readers are closed

	// sort resulting merged values
	for i := range actual {
		slices.Sort(actual[i].postings)
		actual[i].postings = slices.Compact(actual[i].postings)
	}

	require.Equal(t, expected, actual)
}

func TestItReturnsErrorFromInnerMergingIterator(t *testing.T) {
	expectedError := fmt.Errorf("inner error")
	i1 := NewSliceIterator([]TermValues{})
	i2 := NewDynamicSliceIterator(
		func() ([]TermValues, error) {
			return nil, expectedError
		},
		func() error { return nil },
	)

	i3 := NewMergingIterator[TermValues]([]Iterator[TermValues]{i1, i2}, CompareTermValues, MergeTermValues)
	_, err := i3.Next()
	require.ErrorIs(t, err, expectedError)
}
