# Go Pearls

I will be extending this document with some findings that took my attention.
Some of them are pretty trivial, it's expected. Feel free to use these notes for
your study.

#### Manipulating pointers. Reversing a linked list.

In this example I turn a list of integers into a chain of linked nodes and reverse it using 1 temporary pointer.

```go
package main

import (
	"fmt"
	"strings"
	"strconv"
)

type Node struct {
	value int
	next *Node
}

func nodeChain(vals []int) *Node {
	var root, prev, next *Node
	for _, v := range vals {
		next = &Node{v, nil}
		if prev != nil {
			prev.next = next
		}
		prev = next
		if root == nil {
			root = next
		}
	}
	return root
}

func (node *Node) String() string {
	bld := strings.Builder{}
	ptr := node
	for ptr != nil {
		if bld.Len() != 0 {
			bld.WriteString(", ")
		}
		bld.WriteString(strconv.Itoa(ptr.value))
		ptr = ptr.next
	}
	return bld.String()
}

func reverse(node *Node) *Node {
	var head, tail, tmp *Node
	head = node
	for head.next != nil {
		tmp = head.next
		head.next = tail
		tail = head
		head = tmp
	}
	head.next = tail
	return head
}

func main() {
	ch := nodeChain([]int{1, 2, 3, 4, 5, 6, 7})
	fmt.Println(ch)
	rch := reverse(ch)
	fmt.Println(rch)
}
```

#### Extending primitives. A stack example.

In this example I'm extending the basic []int type by adding stack methods to
*stack receiver. 

```go
package main

import (
	"fmt"
)

type stack []int

func (s *stack) Push(v int) {
    *s = append(*s, v)
}

func (s *stack) Peek() int {
    return (*s)[len(*s)-1]
}

func (s *stack) Pop() int {
    v := (*s)[len(*s)-1]
    *s = (*s)[:len(*s)-1]
    return v
}

func (s *stack) IsEmpty() bool {
    return len(*s) == 0
}

func (s *stack) Size() int {
    return len(*s)
}

func main() {
	s := new(stack)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	fmt.Printf("s: %#v, isEmpty: %t\n", s, s.IsEmpty())
	s.Pop()
	fmt.Printf("s: %#v, isEmpty: %t\n", s, s.IsEmpty())
	s.Pop()
	s.Pop()
	fmt.Printf("s: %#v, isEmpty: %t\n", s, s.IsEmpty())
	s.Pop() // panics, as expected
}
```

Note the pointer dereferencing happening all around: it
introduces some vsual noise. We might have defined the same methods for plain stack
type receiver instead of a pointer type, but in this case any modifications to
a stack instance would have to be passed back from the methods explicitly. E.g.

```go
stack = stack.Push(1)
stack, val := stack.Pop()
if stack.IsEmpty() { ... }
```
It's a matter of taste, I assume. The implementation above looks a bit cleaner to me.

#### Binary trees

This paragraph is not Go-specific but a lot of people are missing out the fact
binary trees might be represented using an array and simplify the level of
abstractions needed to implement the basic functionality.


Some facts about binary tree implementations:
* There are 2 major ways of implementing a binary tree: a node struct with 2
  child pointers or an array. The first solution provides a great readability
  and is very easy to understand. The 2nd approach comes with a tiny little bit
  of arithmetics on child/parent index calculation.
* Complete binary trees are very much different from sparse trees in terms of
  amortized performance. The latter ones might be suffering from height
  inflation and therefore perfromance reduction.
* Complete binary tree represented as an array is probably the fastest
  implementation one could achieve. It implies on sequentual memory reading on
  searching Vs random pointer-based lookup chains.
* There is almost no difference between array-based and pointer-based in terms
  of implementing the basic functionality. Both iterative and recursive
  solutions work well and easy with both cases.
* One says using an array for sparse trees is a waste of memory. This is true,
  but dealing with trees with a high level of sparseness is a signal there is a
  better candidate for the data structure in use: RB-trees, AVL-trees, etc.

Example implementation of a binary tree using pointers:

```go
/*
     foo
    /  \
  bar  baz

*/
type node struct {
	value string
	left, right *node
}
foo, bar, baz := &node{"foo"}, &node{"bar"}, &node{"baz"}
foo.left, foo.right = bar, baz

```

The same binary tree represented as an array:

```go
tree := []string{"foo", "bar", "baz"}
```

Array-based trees come with a few tricks, which are worth mentioning:
* if ix is an index of a node, then:
  * it's parent has an index (ix - 1) / 2
  * it's children (left and right accordingly) have indices:
    2 * ix + 1 and 2 * ix + 2
  * it's up to a developer how to indicate an empty placeholder in an
    array-based implementation. On a contrast with pointer-based implementation,
    where a node struct is introduced and a nil-pointer indicates an empty
    child, in this case it might be an extra node struct or a prohibited value
    range (say, an empty string indicates an empty child, or -1 for int values).

An array-based solutiuon comes with some really cool perks: once the index
calculation is rock-solid, the rest might come almost for free.

Here is an exmple of finding the least common ancestor of 2 nodes in a binary
tree using an array-based implementation (assuming, we already know the
indices. If it's not the case and only the lookup values are known in advance,
the tree has to be traversed (O(num elements in the tree) time complexity),
which is very much as expensive as a pointer-based LCA search)

```go
/*
         1
      2     3
   4   5   6  7
 8  9
*/
tree := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
ix1, ix2 := 4, 8 // we are looking for an LCA for nodes 5 and 9
for ix1 != ix2 { // eventually ix1 and ix2 should match on their LCA
	if ix1 > ix2 {
		ix1 = (ix1 - 1) / 2 // set ix1 to it's parent index
	} else {
		ix2 = (ix2 - 1) / 2 // set ix2 to it's parent index
	}
}
// now both ix1 and ix2 contain the LCA index = 1
```

Compare this solution to a pointer-based LCA lookup. In golang it's around 30
lines of code requiring a decent validation.

## Copyrights

Created by Oleg Sidorov in 2019
