package main

import "fmt"

type ListNode struct {
	Next *ListNode
	Prev *ListNode
	Key  int
	Val  int
}

func NewListNode(key int, val int) *ListNode {
	return &ListNode{
		Key: key,
		Val: val,
	}
}

type LRUCache struct {
	head     *ListNode
	tail     *ListNode
	hm       map[int]*ListNode
	length   int
	capacity int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		hm:       make(map[int]*ListNode),
		capacity: capacity,
	}
}

func (lru *LRUCache) Get(key int) int {
	if _, ok := lru.hm[key]; !ok {
		return -1
	}
	node := lru.hm[key]
	lru.promote(node)
	return node.Val
}

func (lru *LRUCache) Put(key int, val int) {
	if _, ok := lru.hm[key]; !ok {
		lru.hm[key] = NewListNode(key, val)
		lru.length++
	}
	node := lru.hm[key]
	node.Val = val
	lru.promote(node)
	for lru.length > lru.capacity {
		lru.trim()
	}
}

func (lru *LRUCache) promote(node *ListNode) {
	if node == lru.head {
		return
	}
	head, tail := lru.head, lru.tail
	prev, next := node.Prev, node.Next
	lru.head = node
	lru.head.Prev = head
	lru.head.Next = nil
	if head != nil {
		head.Next = lru.head
	}
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	if tail == node {
		tail = next
	}
	if tail == nil {
		tail = node
	}
	lru.tail = tail
}

func (lru *LRUCache) trim() {
	if lru.tail == nil {
		return
	}
	k := lru.tail.Key
	delete(lru.hm, k)
	lru.length--
	if lru.head == lru.tail {
		lru.head = nil
		lru.tail = nil
		return
	}
	lru.tail = lru.tail.Next
	if lru.tail != nil {
		lru.tail.Prev = nil
	}
}

func (lru *LRUCache) String() string {
	res := make([]string, 0, lru.length)
	for t := lru.tail; t != nil; t = t.Next {
		k, v := t.Key, t.Val
		kv := fmt.Sprintf("{%d -> %v}", k, v)
		res = append(res, kv)
	}
	return fmt.Sprintf("%+v", res)
}

func main() {
	lru := Constructor(10)
	input := [][]int{{10, 13}, {3, 17}, {6, 11}, {10, 5}, {9, 10}, {13}, {2, 19}, {2}, {3}, {5, 25}, {8}, {9, 22}, {5, 5}, {1, 30}, {11}, {9, 12}, {7}, {5}, {8}, {9}, {4, 30}, {9, 3}, {9}, {10}, {10}, {6, 14}, {3, 1}, {3}, {10, 11}, {8}, {2, 14}, {1}, {5}, {4}, {11, 4}, {12, 24}, {5, 18}, {13}, {7, 23}, {8}, {12}, {3, 27}, {2, 12}, {5}, {2, 9}, {13, 4}, {8, 18}, {1, 7}, {6}, {9, 29}, {8, 21}, {5}, {6, 30}, {1, 12}, {10}, {4, 15}, {7, 22}, {11, 26}, {8, 17}, {9, 29}, {5}, {3, 4}, {11, 30}, {12}, {4, 29}, {3}, {9}, {6}, {3, 4}, {1}, {10}, {3, 29}, {10, 28}, {1, 20}, {11, 13}, {3}, {3, 12}, {3, 8}, {10, 9}, {3, 26}, {8}, {7}, {5}, {13, 17}, {2, 27}, {11, 15}, {12}, {9, 19}, {2, 15}, {3, 16}, {1}, {12, 17}, {9, 1}, {6, 19}, {4}, {5}, {5}, {8, 1}, {11, 7}, {5, 2}, {9, 28}, {1}, {2, 2}, {7, 4}, {4, 22}, {7, 24}, {9, 26}, {13, 28}, {11, 26}}
	for _, i := range input {
		switch len(i) {
		case 1:
			v := lru.Get(i[0])
			fmt.Printf("GET %d -> %d %s\n", i[0], v, lru.String())
		case 2:
			lru.Put(i[0], i[1])
			fmt.Printf("PUT %d -> %d %s\n", i[0], i[1], lru.String())
		default:
			panic("unexpected input")
		}
	}
}
