# go-wbbst

A golang implementation of the Weight Balanced Binary Search tree implemented without vectors.

## General Description of Algorithm

Binary Search Trees are usually implemented using a 4 part data structure, and as such they generally use around 4 times as much memory as would be required for a data payload the size of the pointers used in them. There is the parent, the left and the right link, and then the data itself.

Weight Balanced Binary Search Trees are a subtype of binary search tree that optimises the number of nodes on each side of a parent or root node and are commonly used in Functional programming, mainly because this particular search tree strategy is very suited to recursion.

This project is an implementation of a Weight Balanced Binary Search Tree in Golang that takes an entirely different approach to the tree structure. By their nature, WBBSTs tend to grow new rows more slowly than other types of BSTs because they try to fill the tree from the top down.

While I was working on a partial hash collision search algorithm for a variant of Cuckoo Cycle, as I started to implement the search strategy for it, using a WBST, I realised that an entirely different approach could be taken to storing the data in the nodes that holds the structure implicitly instead of explicitly using pointers.

go-wbbst stores the tree structure instead as an array of the data payload data type, in the case of my project it was 32 and 64 bit words, in a novel structure that I am not aware of being used previous to this invention.

A binary search tree, if you consider it as a collection of empty slots for your data, has 1, 2, 4, 8, 16, and so on, with each row potentially storing double the number of data objects as the previous.

This is mainly applicable to the case of storing an index in a form that can be rapidly searched, items inserted, and deleted, with an extremely fast and compact representation that especially speeds up search and insert. For deletion it is somewhat slower, as the larger the tree gets, the more data has to be shifted upwards towards the root.

The application that it was dreamed up for consisted of progressively generating hashes, splitting them, and searching for half collisions. However, this can be generalised to allow the storage of any type of data that can be easily evaluated as greater than or less than. Or in other words, it is a fast index algorithm that can be used for data that is searched atomically. Input could be hashed, and then extremely rapidly the tree can be walked to look for a match, or a place to slot the hash in, when growing the index, and less quickly, deleted from the tree. It is really most applicable to a search index that will live largely in memory, and with on-die CPU caches over 4mb, they will live largely in cache and thus give the extra advantage of the extreme low latency and high throughput of this memory.

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

So the right of 4 is 2\*4+2=10, the left of 5 is 2\*5+1=11. Since doubling and halving can be performed entirely with bitshifting, finding which index to make a walk is 1 left rotation and 1 addition, or walking the other way, a subtraction or addition and right rotation, depending on whether the number is even or odd (this can be discovered also by a simple bitwise AND operation with 1).

## Tree Rotations

Because this implements a Weight Balanced Binary Search Tree, for any insert operation there can be one or two rotations required.

The changes required can be determined before allocating a subsequent row to the dynamic array slices, as we know from a comparison to the root 0 data if we are going to the left or right, the rotation moves a new lesser child from two rows up to the same side of the parent as it is a child of its parent, the displaced object moves to become the same side child of this new parent.

Here is an example:

|     |     |     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 10  |     |     |     |     |     |     |     |
| 12  | 8   |     |     |     |     |     |     |
| 14  | 11  |     | 4   |     |     |     |     |
| 16  |     |     |     |     |     | 7   | 2   |

This is an imbalanced tree. 8 does not have a left node. To balance it, 7 replaces 8, and 8 goes in beside the 4 as the left of 8

|     |       |       |     |     |     |     |     |
| --- | ----- | ----- | --- | --- | --- | --- | --- |
| 10  |       |       |     |     |     |     |     |
| 12  | **7** |       |     |     |     |     |     |
| 14  | 11    | **8** | 4   |     |     |     |     |
| 16  |       |       |     |     |     |     | 2   |

You can quickly determine the tree is imbalanced in the first example because of the empty space between 11 and 4, and the 7 and 2 on the right side, which are children of 4, whereas 8 only has 4 as a child.

The objective a Weight Balanced tree is to produce the minimum number of steps between the root and any other node in the tree. The top example has three paths that require 3 hops, whereas the second one has only two.
