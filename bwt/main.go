package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		panic("Unable to read form STDIN")
	}
	str = strings.TrimRight(str, "\r\n")
	rotations := make([]string, 0, len(str))
	rotations = append(rotations, str)
	for i := 1; i < len(str); i++ {
		rotations = append(rotations, rotate(str, i))
	}
	sort.Strings(rotations)
	writer := strings.Builder{}
	for _, s := range rotations {
		writer.WriteByte(s[len(s)-1])
	}

	fmt.Println(writer.String())
}

func rotate(str string, shift int) string {
	length := len(str)
	return str[length-shift:length] + str[:length-shift]
}
