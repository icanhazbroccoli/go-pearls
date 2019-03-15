package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func allIndicesOf(str, pattern string) []int {
	s := pattern + "\x00" + str
	res := make([]int, 0)

	if len(pattern) > len(str) {
		return res
	}

	prefix := make([]int, len(s))
	var border int
	for i := 1; i < len(s); i++ {
		for border > 0 && s[border] != s[i] {
			border = prefix[border-1]
		}
		if s[i] == s[border] {
			border = border + 1
		} else {
			border = 0
		}
		prefix[i] = border
	}

	l := len(pattern)
	for ix, border := range prefix {
		if l == border {
			res = append(res, ix-2*l)
		}
	}

	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	pattern, err := reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	pattern = strings.TrimRight(pattern, "\r\n")

	str, err := reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	str = strings.TrimRight(str, "\r\n")

	res := allIndicesOf(str, pattern)

	fmt.Println(res)
}
