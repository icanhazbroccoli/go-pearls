# Go Pearls

I will be extending this document with some findings that took my attention.
Some of them are pretty trivial, it's expected. Feel free to use these notes for
your study.

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
