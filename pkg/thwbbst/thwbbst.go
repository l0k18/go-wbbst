// Package thwbbst - A Top Heavy Weight Balanced Binary Search Tree for sorting small, largely random blobs of data
package thwbbst

import "errors"

// Wbbst - Top Heavy Weight Balanced Binary Search Tree
// Tree is encoded by an array with uint32 indices. Data type is unspecified and defined by the implementation.
type Wbbst interface {
	// Comparator functions. Name indicates truth value from first argument to the second (IsLeft is true if first arg is left of second)
	// These are an interface to allow implementation different data blob types. For example a 64 bit hash, but only compare the first 32.
	// It could also be used for other types of compact, complex values like vectors and bit matrixes
	IsLeft(interface{}, interface{}) bool
	IsRight(interface{}, interface{}) bool
	// Returns the index location of a piece of data or an error if not found
	Find(interface{}) (uint32, error)
	// Adds a new row to the bottom of the tree (or error if row would extend beyond 32bits of index). This must increment Depth and make and append a string of array elements the length of the 2 to the power of the new Depth (1<<Depth)
	AddRow() error
	// Insert a new node in the tree and balance if necessary. Returns position of insertion.
	Insert(interface{}) uint32
	// Removes item at index and rebalances. Presumably zero is the sentinel for no allocation. The sentinal is entirely implementation driven but
	// by default golang zeroes all new variables it creates to avoid potential security vulnerabilities this can introduce, as well as eliminating
	// 'may be used uninitialised' bugs. If other than zero must be used as a sentinel, AddRow() must be reimplemented also to suit.
	DeleteByIndex(uint32)
	// Searches for a node with specified data and rebalances.
	DeleteByData(interface{}) error
	// Walking functions. These take an index and return the correct index of the data on this path. These functions are the only ones this package actually implements, as everything else depends on the comparators
	WalkUp(uint32) (uint32, error)
	WalkLeft(uint32) (uint32, error)
	WalkRight(uint32) (uint32, error)
}

// Tree - a generic structure for storing a short data type for the search tree
type Tree struct {
	Weight [2]uint32
	Depth  uint8
	// The store is tree structured only by the inherent structure of the tree (it is just rows of progressively doubling strings)
	Store []interface{}
}

// AddRow - add a new row to the bottom of the tree
func (t *Tree) AddRow() error {
	if 1<<t.Depth == 0 {
		return errors.New("Tree is already at maximum depth of 31")
	}
	t.Store = append(t.Store, make([]interface{}, 1<<t.Depth))
	t.Depth++
	return nil
}

// WalkUp - return the correct index from the Store that corresponds to the Parent of the index argument
func (t *Tree) WalkUp(index uint32) (uint32, error) {
	if index == 0 {
		return 0, errors.New("Cannot walk up from the root")
	}
	return index>>1 - (index+1)%2, nil
}

// WalkLeft - Return the index from the Store that corresponds to the Left value in the next row down
func (t *Tree) WalkLeft(index uint32) (uint32, error) {
	left := index << 1
	if left > 1<<t.Depth {
		return 0, errors.New("Cannot walk below the bottom of the tree")
	}
	return left + 2, nil
}

// WalkRight - Return the index from the Store that corresponds to the Right value in the next row down
func (t *Tree) WalkRight(index uint32) (uint32, error) {
	right := index << 1
	if right > 1<<t.Depth {
		return 0, errors.New("Cannot walk below the bottom of the tree")
	}
	return right + 2, nil
}
