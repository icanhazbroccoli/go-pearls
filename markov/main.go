package main

import (
	"bufio"
	"os"
)

const (
	MaxGen  = 10000
	Npref   = 2
	NonWord = ""
)

type Prefix struct {
	length int
	items  [Npref]string
}

func (p *Prefix) PushBack(s string) {
	complete := p.IsComplete()
	ix := p.length
	if p.length == Npref {
		for i := 1; i < Npref; i++ {
			p.items[i-1] = p.items[i]
		}
	} else {
		p.length++
	}
	p.items[ix] = s
}

func (p *Prefix) IsComplete() bool {
	return p.length == Npref
}

var (
	statetab map[Prefix][]string
)

func init() {
	statetab = make(map[Prefix][]string)
}

func build(prefix *Prefix, scanner *bufio.Scanner) {
	for scanner.Scan() {
		add(prefix, scanner.Text())
	}
}

func add(prefix *Prefix, s string) {
	if prefix.length == Npref {
		if _, ok := statetab[*prefix]; !ok {
			statetab[*prefix] = make([]string, 0, 1)
		}
		statetab[*prefix] = append(statetab[*prefix], s)
		for i := 1; i < Npref; i++ {
			prefix.items[i-1] = prefix.items[i]
		}
	}
	prefix.items[prefix.length] = s
	if prefix.length < Npref {
		prefix.length++
	}
}

func generate(nwords int, w *bufio.Writer) {
	prefix := &Prefix{}
	for i := 0; i < nwords; i++ {

	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	prefix := &Prefix{}
	build(prefix, scanner)
}
