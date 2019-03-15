package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Base = 65521
	Offs = 1
)

func adler32(s string) uint32 {
	var a, b uint32 = 1, 0
	for _, ch := range s {
		a = (a + uint32(ch)) % Base
		b = (b + a) % Base
	}
	return (b << 16) | a
}

func rotateRight(adler uint32, length int, left, right byte) uint32 {
	var a, b uint32 = adler & 0xFFFF, adler >> 16
	a = (a - uint32(left) + uint32(right)) % Base
	b = (b - uint32(length)*uint32(left) + a - Offs) % Base

	return (b << 16) | a
}

func rksearch(text string, pattern string) []int {
	fmt.Printf("Text: %q, pattern: %q\n", text, pattern)
	res := make([]int, 0)
	if len(pattern) > len(text) {
		return res
	}
	phash := adler32(pattern)
	plen := len(pattern)
	tlen := len(text)

	for i, thash := 0, adler32(text[0:plen]); i < tlen-plen+1; i++ {
		fmt.Printf("%q [%d] Vs %q [%d]\n", text[i:i+plen], thash, pattern, phash)
		if phash == thash {
			fmt.Println("Hash match")
			if pattern == text[i:i+plen] {
				fmt.Printf("Match on position %d\n", i)
				res = append(res, i)
			}
		}
		if i == tlen-plen {
			break
		}
		fmt.Printf("-%c, +%c\n", text[i], text[i+plen])
		thash = rotateRight(thash, plen, text[i], text[i+plen])
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	text = strings.TrimRight(text, "\r\n")

	nstr, err := reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	n, err := strconv.Atoi(strings.TrimRight(nstr, "\r\n"))
	if err != nil {
		panic(err.Error())
	}

	if n <= 0 {
		panic(fmt.Sprintf("n is expected to be greater than 0, %d provided", n))
	}

	patterns := make([]string, 0, n)
	for i := 0; i < n; i++ {
		pattern, err := reader.ReadString('\n')
		if err != nil {
			panic(err.Error())
		}
		patterns = append(patterns, strings.TrimRight(pattern, "\r\n"))
	}

	for _, pattern := range patterns {
		fmt.Println(rksearch(text, pattern))
	}
}
