# go-wbbst

A golang implementation of the Weight Balanced Binary Search tree implemented without vectors.

## Contents

* [How to use go-wbbst](#how-to-use-go-wbbst)
* [General Description of Algorithm](#general-description-of-algorithm)
* [Tree Walk Functions](#tree-walk-functions)
* [Data Structure](#data-structure)
* [Tree Rotations](#tree-rotations)
* [Memory Architecture Considerations - Why this will be faster](#memory-architecture-considerations---why-this-will-be-faster)
* [What inspired go-wbbst?](#what-inspired-go-wbbst)

## How to use go-wbbst

This library only provides the walk functions and a skeleton data structure that you will want your functions to use.

In order to use it you must write your implementation of comparators, search, insert and delete functions (depending on what your application requires).

## General Description of Algorithm

Binary Search Trees are usually implemented using a 4 part data structure, and as such they generally use around 4 times as much memory as would be required for a data payload the size of the pointers used in them. There is the parent, the left and the right link, and then the data itself.

Weight Balanced Binary Search Trees are a subtype of binary search tree that optimises the number of nodes on each side of a parent or root node and are commonly used in Functional programming, mainly because this particular search tree strategy is very suited to recursion.

This project is an implementation of a Weight Balanced Binary Search Tree in Golang that takes an entirely different approach to the tree structure. By their nature, WBBSTs tend to grow new rows more slowly than other types of BSTs because they try to fill the tree from the top down.

While I was working on a partial hash collision search algorithm for a variant of Cuckoo Cycle, as I started to implement the search strategy for it, using a WBST, I realised that an entirely different approach could be taken to storing the data in the nodes that holds the structure implicitly instead of explicitly using pointers.

go-wbbst stores the tree structure instead as an array of the data payload data type, in the case of my project it was 32 and 64 bit words, in a novel structure that I am not aware of being used previous to this invention.

A binary search tree, if you consider it as a collection of empty slots for your data, has 1, 2, 4, 8, 16, and so on, with each row potentially storing double the number of data objects as the previous.

This is mainly applicable to the case of storing an index in a form that can be rapidly searched, items inserted, and deleted, with an extremely fast and compact representation that especially speeds up search and insert. For deletion it is somewhat slower, as the larger the tree gets, the more data has to be shifted upwards towards the root.

The application that it was dreamed up for consisted of progressively generating hashes, splitting them, and searching for half collisions. However, this can be generalised to allow the storage of any type of data that can be easily evaluated as greater than or less than. Or in other words, it is a fast index algorithm that can be used for data that is searched atomically. Input could be hashed, and then extremely rapidly the tree can be walked to look for a match, or a place to slot the hash in, when growing the index, and less quickly, deleted from the tree. It is really most applicable to a search index that will live largely in memory, and with on-die CPU caches over 4mb, they will live largely in cache and thus give the extra advantage of the extreme low latency and high throughput of this memory.

## Tree Walk Functions

Instead of using references or pointers, the array indices themselves, when used in the correct formula, provide the correct index for the desired walk. As you will see below, if implemented in Assembly Language the up-walk takes 4 instructions, and the left/right walks only require two instructions.

They do need to be prefixed by a test. For up, if the index is zero, you can't go up from there. For left and right, the depth of the array - 2 to the power of the number of rows, is the limit. Thus in fact the up walk takes 5 instructions and the left/right walks take 4 (bitshift, conditional branch, bitshift, addition)

i = index of array element

Walk Up (go to parent):

`i>>1-(i+1)%2`

(Obviously you probably need to test if i is zero)

Walk Left (down and left):

`i<<1+2`

(this and the next would also need to test against 2^depth of the structure, so it doesn't walk off the edge and fall back to the root)

Walk Right (down and right):

`i<<1+1`

## Data Structure

The data structure consists of an array, implemented in Golang as a dynamic slice array, and each row exists in progressively doubling numbers of indices of the array. Thus row 1 has only index zero, row two has 1, 2, row 3 has 4, 5, 6, 7, row 4 has 8, 9, 10, 11, 12, 13, 14, 15 and so on.

When walking the tree, using greater/less comparisons, there is a simple formula that determines the index that gives you the left and right node under the parent, for example, the left node from 2 is 6, the right node of 6 is 13.

Visually, it looks like this, obviously it is not practical to illustrate it with much more width than 16 columns wide.

|     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 0   |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| 1   | 2   |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| 3   | 4   | 5   | 6   |     |     |     |     |     |     |     |     |     |     |     |     |
| 7   | 8   | 9   | 10  | 11  | 12  | 13  | 14  |     |     |     |     |     |     |     |     |
| 15  | 16  | 17  | 18  | 19  | 20  | 21  | 22  | 23  | 24  | 25  | 26  | 27  | 28  | 29  | 30  |

The left and right of 1 is 3, 4, the left and right of 13 is 27, 28, the left and right of 10. You can find the position of left and right by doubling the number and adding 1 for left or 2 for right.

So the right of 4 is 2\*4+2=10, the left of 5 is 2\*5+1=11. Since doubling and halving can be performed entirely with bitshifting, finding which index to make a walk is 1 left rotation and 1 addition, or walking the other way, a subtraction or addition and right rotation, depending on whether the number is even or odd (this can be discovered also by a simple bitwise AND operation with 1). Or, as I have specified in the foregoing, using an addition to the index and a modulus.

## Tree Rotations

Because this implements a Weight Balanced Binary Search Tree, for any insert operation there can be one or two rotations required.

The changes required can be determined before allocating a subsequent row to the dynamic array slices, as we know from a comparison to the root 0 data if we are going to the left or right, the rotation moves a new lesser child from two rows up to the same side of the parent as it is a child of its parent, the displaced object moves to become the same side child of this new parent.

Here is an example of a double rotation:

```
                           8
               12                      4
        14           11          7            2
    16      --    --    --   --     --    --     --
```

We want to add the number 10 to this structure. Starting from the root we go left, hit 12, so we go right, and we hit 11, so we go right and find an empty slot.

```
                           8
               12                      4
        14           11          7            2
    16      --   --      10  --     --    --     --
```

However, now we have two 3 step paths to the left, and none to the right. A weight balanced tree would have one such on each side.

Thus we must recenter on the lowest value from the left side, our new 10. The 10 moves up, and pushes the 8 to its right, and then we have another imbalanced tree:

```
                           10
               12                      8
        14           11          --           4
    16      --   --      --  --      --   7       2
```

Now it's even worse, we have three three step paths, though it is balanced left to right. For this reason I also considered to call this algorithm 'Top Heavy Weight Balanced Binary Search Tree', and this is why this insert required two rotations.

To balance this, we need to rotate the 7 up to the position of the 8, the 8 goes to the left of the 7:

```
                           10
               12                       7
        14           11           8           4
    16      --   --      --    --   --     --    2
```

(This wing-shaped result will not always be the rule, but will happen a lot, and has a poetic symmetry considering the avian theme of my other work)

Now the tree is balanced, and we have two 3 step paths, on each side, and every other slot in the rows is only 2 step paths, and every spot is filled.

Now, instead of even putting the 10 in and doing all of these rotations, we can determine straight away that as soon as we saw the 8 at the root, and we travelled left to insert the 10, we know because there was already 4 nodes left and 3 nodes right, that we have to perform a rotation. Rotations that shift the root node are by necessity double rotations because of that hole that appears in the first step.

So the algorithm instead then knows we have to recenter on the 10, as it is to the right of the rightmost left leaf, that then the 8 moves to the right of the root, but we know because the left-most leaf of the right hand side is to the right of 8, that we should put the 8 to the left of next right-most, the 7 instead, becomes the right of the root, and the 4 also had to shift down to the right, and the 2 down to the right of it.

If the new node goes to the side already heavier than the right, then we know the root has to move right, and all other right-children have to also shift to the right. The transformation is recursive.

The big benefit of using this bifurcating linear array with doubling length substrings means that when this rotation has to occur we know we have to recenter immediately from the first comparison. We know that the right hand edge of the tree must shift to the right (and thus down) one layer, and in doing this, the old right of the root must go to the left of the top node we shifted right.

To perform this most simply, we then follow the right hand edge down until we hit an empty row, copy the parent into the leaf, then move up, copy the parent, until the next is the root, we place our new root into the root, and we move the old root to the left of the old root.

## Memory Architecture Considerations - Why this will be faster

Of course in a very large tree, this operation could span perhaps even many scores of rows cascading downwards. Golang only has 32 bit array indices, so if our data was 64 bits, this array could get no deeper, if every value did not repeat, than 32.

However, compared to a conventional binary tree, with 32 bit indices (or pointers) from each node, then we have 20 bytes to store just one node. This would mean 72 gigabytes of storage for this theoretical comprehensive tree.

With this structure instead we are looking at merely 16gb of storage, with a small overhead for tracking left right balance and the depth of the tree.

Furthermore, you have to consider the way that the CPU will cache the data. With a conventional 4 part data structure and bidirectional references on each node, 20 bytes have to be loaded into the CPU cache to examine it. 20 bytes is equal to 5 nodes with Top-Heavy WBBST.

If you consider a typical mid-range CPU of current vintage, around 6Mb of cache, if the CPU was doing nothing else, it could store 786,232 nodes.

This would mean 16 full rows starting from the root, and it would use about 2/3 of the available storage, so this search structure could stay largely in cache during these rotation operations, and then the altered data would be written back to memory in a string of big dumps. This means this algorithm will tend mostly towards a large amount of linear memory transfers, in comparison to the amount of tree structure we can examine.

## What inspired go-wbbst?

This extreme imbalance will be fully a ratio of 1:4 for a hash table of 32 bit long hashes. This was the payload that I was considering when I came up with this algorithm, and I didn't want to implement the search in such an inefficient way, I wanted the solver this is part of, to be very hard to optimise in any way to improve performance. CPUs are way bottlenecked for memory access compared to GPUs, but their on-die caches are massively bigger, and way faster than any kind of GPU memory.

This search data structure will naturally shift the balance back to the CPU as being faster at dealing with an Extreme-Memory-Hard Proof of Work, because of its cache, and more than doubly so in the case of the Ryzen 7 CPU with 20Mb of onboard cache. Similar to the Cryptonight algorithm, except not limited to merely 2gb per thread for performing fast hashes in on-die cache, this will benefit from as much cache as you can throw at it, further reducing random access of memory and thus skipping past the far greater retrieval latency of DDR4 memory compared to GDDR5X, which will then be bottlenecked by its lower latency, where the CPU can handle large chunks of this tree structure without having to make nonlinear requests from memory.

Beyond the immediate problem this data structure is intended to solve, this binary tree search library can also massively accelerate a database index, hence why I have written it into a library with a generalised untyped array interface so it can be adapted to work with any readily sortable hash table search.
