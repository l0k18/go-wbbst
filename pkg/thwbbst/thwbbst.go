// Package thwbbst - A Top Heavy Weight Balanced Binary Search Tree for sorting small, largely random blobs of data
package thwbbst

// Wbbst - Top Heavy Weight Balanced Binary Search Tree
// Tree is encoded by an array with uint32 indices. Data type is unspecified and defined by the implementation.
// This interface can be implemented with any alternative type of BST implementation but this compact storage format is very well suited to scapegoat and weight balanced, and the extra data in the search data structure will be required for each implementation. See the tree32 implementation in /pkg in this repository
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
