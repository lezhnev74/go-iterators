# Go Iterators

[![Go package](https://github.com/lezhnev74/go-iterators/actions/workflows/go.yml/badge.svg)](https://github.com/lezhnev74/go-iterators/actions/workflows/go.yml)
![Code Coverage](https://raw.githubusercontent.com/lezhnev74/go-iterators/badges/.badges/main/coverage.svg)

Since Go does not have a default iterator type (though there are
discussions [here](https://bitfieldconsulting.com/golang/iterators), [here](https://github.com/golang/go/issues/61897)
and [there](https://ewencp.org/blog/golang-iterators/)), here is a set of different iterators crafted manually.
Particularly, there is [a proposal](https://github.com/golang/go/issues/61898) for a package that defines compound
operations on iterators, like merging/selecting. Until Go has a stdlib's iterator implementation (or at least an
experimental standalone package), there is this package.

## Iterator Interface

```go
// Iterator is used for working with sequences of possibly unknown size
// Interface adds a performance penalty for indirection.
type Iterator[T any] interface {
// Next returns EmptyIterator when no value available at the source
// error == nil means the returned value is good
Next() (T, error)
// The client may decide to stop the iteration before EmptyIterator recieved
// Closed iterator may panic
io.Closer
}
```

## Various Iterators

- `CallbackIterator` calls a callback to fetch the next value
- `SliceIterator` iterates over a static precalculated slice
- `DynamicSliceIterator` behaves like `SliceIterator` but fetches a new slice when previous slice is "empty"
- Selection-tree iterators (only for sorted iterators). Selecting iterators are to form a selection tree that helps
  efficiently merge sorted values with the least number of compare operations.
    - `SelectingIterator` combines 2 sorted iterators. Effectively that is a set union.
    - `UniqueSelectingIterator` The same as `SelectingIterator` but removes duplicates.
    - `RemovingIterator` combines 2 sorted iterators and removes values from one that is present in the second.
      Effectively that is a set difference.
    - `SortedSelectingIterator` combines 2 sorted iterators into a single sorted iterator. 

