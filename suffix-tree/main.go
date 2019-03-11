package main

import "fmt"

type TreeNode struct {
	start    uint32
	length   uint32
	children map[uint64]*TreeNode
}

func NewTreeNode(start, length uint32) *TreeNode {
	return &TreeNode{
		start:    start,
		length:   length,
		children: make(map[uint64]*TreeNode),
	}
}

type SuffixTree struct {
	text string
	root *TreeNode
}

func NewSuffixTree(text string) *SuffixTree {
	t := &SuffixTree{
		text: text,
		root: NewTreeNode(0, 0),
	}
	t.index()

	return t
}

func (t *SuffixTree) index() {
	var branch *TreeNode
	var branchstartlength uint64
	for i := 0; i < len(t.text); i++ {
		ptr := t.root
		off := uint32(i)

		for {
			branch = nil
			branchstartlength = 0
			for startlength, ch := range ptr.children {
				start, _ := decodestartlength(startlength)
				if start >= uint32(len(t.text)) {
					continue
				}
				if t.text[off] == t.text[start] {
					// first character match, this is the
					// right branch
					branch = ch
					branchstartlength = startlength
					break
				}
				continue
			}
			if branch == nil {
				startlength := encodestartlength(off, uint32(len(t.text))-off)
				ptr.children[startlength] = NewTreeNode(off, uint32(len(t.text))-off)
				break
			}
			branchstart, branchlength := decodestartlength(branchstartlength)
			diverged := -1
			// 1 because we have already found the match
			for j := uint32(1); j < branchlength; j++ {
				if off+j >= uint32(len(t.text)) || t.text[off+j] != t.text[branchstart+j] {
					diverged = int(j)
					break
				}
			}
			if diverged == -1 {
				off += branchlength
				ptr = branch
				continue
			}
			delete(ptr.children, branchstartlength)
			newnode := NewTreeNode(branchstart, uint32(diverged))
			ptr.children[encodestartlength(branchstart, uint32(diverged))] = newnode
			newnode.children[encodestartlength(off+uint32(diverged), uint32(len(t.text))-(off+uint32(diverged)))] = NewTreeNode(off+uint32(diverged), uint32(len(t.text))-(off+uint32(diverged)))
			branch.start = branchstart + uint32(diverged)
			branch.length = branchlength - uint32(diverged)
			newnode.children[encodestartlength(branch.start, branch.length)] = branch
			break
		}
	}
}

func decodestartlength(startlength uint64) (uint32, uint32) {
	return uint32(startlength >> 32), uint32(startlength & 0xFFFFFFFF)
}

func encodestartlength(start, length uint32) uint64 {
	return uint64(start)<<32 | uint64(length)
}

func main() {
	//text := "banana"
	text := "ATAAATG"
	t := NewSuffixTree(text)
	fmt.Println(t)
}
