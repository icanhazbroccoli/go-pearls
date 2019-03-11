package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func InverseBwt(str string) string {
	right := []rune(str)
	left := make([]rune, len(right))
	copy(left, right)
	sort.Slice(left, func(i, j int) bool {
		return left[i] < left[j]
	})

	mapping := make(map[rune][]int)
	for ix, ch := range right {
		if _, ok := mapping[ch]; !ok {
			mapping[ch] = make([]int, 0, 1)
		}
		mapping[ch] = append(mapping[ch], ix)
	}
	shifts := make([]int, len(left))
	chars := make(map[rune]int)
	for ix, ch := range left {
		shifts[ix] = chars[ch]
		chars[ch]++
	}

	writer := strings.Builder{}
	next := mapping['$'][0]
	for writer.Len() < len(str) {
		nextCh := left[next]
		writer.WriteRune(nextCh)
		next = mapping[nextCh][shifts[next]]
	}

	return writer.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		panic("Unable to read from STDIN")
	}
	str = strings.TrimRight(str, "\r\n")

	orig := InverseBwt(str)

	fmt.Println(orig)
}
