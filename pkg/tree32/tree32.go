// Package tree32 - Implements a Top Heavy Weight Balanced Binary Tree for 32 bit unsigned integers.
package tree32

import "errors"

// Tree - contains a WBBST for uint32 payload
type Tree struct {
	Weight [2]uint32
	Depth  uint32
	Store  []uint32
}

// NewTree - creates a new Tree object
func NewTree() Tree {
	var tree Tree
	tree.Store = make([]uint32, 1)
	tree.Depth++
	return tree
}

// AddRow - adds a new row to the tree
func (t *Tree) AddRow() error {
	oldLen := len(t.Store)
	s := make([]uint32, 1<<t.Depth)
	t.Store = append(t.Store, s...)
	if len(t.Store) == oldLen {
		return errors.New("Unable to add new row to tree")
	}
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
	return right + 1, nil
}

// IsLeft - implements a greater than test of a value on a tree node
func (t *Tree) IsLeft(object, subject uint32) bool {
	return object > t.Store[subject]
}

// IsRight - implements a less than test of a value on a tree node
func (t *Tree) IsRight(object, subject uint32) bool {
	return object < t.Store[subject]
}
