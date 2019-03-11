package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func betterbwmatch(input string, patterns []string) []int {

	charIx := 0
	charMap := make(map[rune]int)
	for _, ch := range input {
		if _, ok := charMap[ch]; !ok {
			charMap[ch] = charIx
			charIx++
		}
	}
	count := make([][]int, len(input)+1)
	count[0] = make([]int, len(charMap))
	for ix, ch := range input {
		count[ix+1] = make([]int, len(charMap))
		copy(count[ix+1], count[ix])
		count[ix+1][charMap[ch]]++
	}

	last := []rune(input)
	first := make([]rune, len(last))
	copy(first, last)
	sort.Slice(first, func(a, b int) bool { return first[a] < first[b] })

	firstOccurence := make(map[rune]int)
	var prevCh rune
	for ix, ch := range first {
		if ch != prevCh {
			firstOccurence[ch] = ix
			prevCh = ch
		}
	}

	res := make([]int, 0, len(patterns))

Pattern:
	for _, pattern := range patterns {
		if len(pattern) == 0 || len(pattern) > len(input) {
			res = append(res, 0)
			continue
		}
		ptr := len(pattern) - 1
		top, bottom := 0, len(input)-1
		for top <= bottom {
			if ptr < 0 {
				res = append(res, bottom-top+1)
				continue Pattern
			}
			ch := rune(pattern[ptr])
			ptr--
			top = firstOccurence[ch] + count[top][charMap[ch]]
			bottom = firstOccurence[ch] + count[bottom+1][charMap[ch]] - 1
			if top > bottom {
				res = append(res, 0)
				continue Pattern
			}
		}
	}

	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	input = strings.TrimRight(input, "\r\n")

	_, err = reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	patternline, err := reader.ReadString('\n')
	patterns := strings.Split(strings.TrimRight(patternline, "\r\n"), " ")

	res := betterbwmatch(input, patterns)

	writer := strings.Builder{}
	for _, pos := range res {
		if writer.Len() > 0 {
			writer.WriteRune(' ')
		}
		writer.WriteString(strconv.Itoa(pos))
	}

	fmt.Println(writer.String())
}
