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
It's a matter of taste, I assume. The variant above looks a bit cleaner to me.

## Copyrights

Created by Oleg Sidorov in 2019
