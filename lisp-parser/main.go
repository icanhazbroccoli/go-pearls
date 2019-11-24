package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Scope map[string]int

func NewScope() Scope {
	return Scope(make(map[string]int))
}

type Node interface {
	Eval(...Scope) (int, error)
	Append(Node) error
}

/* IntNode */
type IntNode struct {
	raw string
}

var _ Node = (*IntNode)(nil)

func NewIntNode(raw string) *IntNode {
	return &IntNode{raw: raw}
}

func (i *IntNode) Eval(...Scope) (int, error) {
	return strconv.Atoi(i.raw)
}

func (i *IntNode) Append(n Node) error {
	return fmt.Errorf("IntNode %q can not append a child node %+v", i.raw, n)
}

/* End of IntNode */

/* VarNode */

type VarNode struct {
	raw string
}

var _ Node = (*VarNode)(nil)

func NewVarNode(raw string) *VarNode {
	return &VarNode{raw: raw}
}

func (v *VarNode) Eval(ss ...Scope) (int, error) {
	name := v.Name()
	for i := len(ss) - 1; i >= 0; i-- {
		if v, ok := ss[i][name]; ok {
			return v, nil
		}
	}
	return 0, fmt.Errorf("variable %q is undefined", name)
}

func (v *VarNode) Append(n Node) error {
	return fmt.Errorf("VarNode %q can not append a child node %+v", v.Name(), n)
}

func (v *VarNode) Name() string {
	return v.raw
}

/* End of VarNode */

/* AddNode */

type AddNode struct {
	nodes []Node
}

var _ Node = (*AddNode)(nil)

func NewAddNode() Node {
	return &AddNode{
		nodes: make([]Node, 0, 2),
	}
}

func (a *AddNode) Eval(ss ...Scope) (int, error) {
	if len(a.nodes) != 2 {
		return 0, fmt.Errorf("AddNode must have exactly 2 operands, got: %d", len(a.nodes))
	}
	subres := make([]int, 0, 2)
	for _, node := range a.nodes {
		r, err := node.Eval(ss...)
		if err != nil {
			return 0, err
		}
		subres = append(subres, r)
	}
	return subres[0] + subres[1], nil
}

func (a *AddNode) Append(n Node) error {
	if len(a.nodes) >= 2 {
		return fmt.Errorf("AddNode can't append more child nodes")
	}
	a.nodes = append(a.nodes, n)
	return nil
}

/* End of AddNode */

/* MultNode */

type MultNode struct {
	nodes []Node
}

func NewMultNode() Node {
	return &MultNode{
		nodes: make([]Node, 0, 2),
	}
}

func (m *MultNode) Eval(ss ...Scope) (int, error) {
	if len(m.nodes) != 2 {
		return 0, fmt.Errorf("MultNode must have exactly 2 operands, got: %d", len(m.nodes))
	}
	subres := make([]int, 0, 2)
	for _, node := range m.nodes {
		r, err := node.Eval(ss...)
		if err != nil {
			return 0, err
		}
		subres = append(subres, r)
	}
	return subres[0] * subres[1], nil
}

func (m *MultNode) Append(n Node) error {
	if len(m.nodes) >= 2 {
		return fmt.Errorf("MultNode can't append more child nodes")
	}
	m.nodes = append(m.nodes, n)
	return nil
}

/* End of MultNode */

/* LetNode */

type LetNode struct {
	nodes []Node
}

var _ Node = (*LetNode)(nil)

func NewLetNode() Node {
	return &LetNode{
		nodes: make([]Node, 0, 3), // at least 1 k->v assignment and an expr
	}
}

func (l *LetNode) Eval(ss ...Scope) (int, error) {
	nodelen := len(l.nodes)
	if nodelen < 3 {
		return 0, fmt.Errorf("LetNode: at least 3 operands expected, %d provided", nodelen)
	}
	if nodelen%2 == 0 {
		return 0, fmt.Errorf("LetNode: even number of operands expected, odd provided (%d)", nodelen)
	}
	i := 0
	scope := NewScope()
	extss := append(ss, scope)
	for i < nodelen-1 {
		namenode, valnode := l.nodes[i], l.nodes[i+1]
		varnode, ok := namenode.(*VarNode)
		if !ok {
			return 0, fmt.Errorf("LetNode: operand %+v is not a var node", namenode)
		}
		val, err := valnode.Eval(extss...)
		if err != nil {
			return 0, err
		}
		scope[varnode.Name()] = val
		i += 2
	}
	return l.nodes[i].Eval(extss...)
}

func (l *LetNode) Append(n Node) error {
	l.nodes = append(l.nodes, n)
	return nil
}

/* End of LetNode */

type parserState uint16

const (
	wantOpeningParantheses parserState = iota
	wantOptionalSpace
	wantAlpha
	wantNumeric
	wantAlphaNumeric
)

var nodeBuilders map[string]func() Node

func init() {
	nodeBuilders = map[string]func() Node{
		"add":  NewAddNode,
		"mult": NewMultNode,
		"let":  NewLetNode,
	}
}

func isSpace(ch byte) bool {
	return ch == ' '
}

func isOpeningParantheses(ch byte) bool {
	return ch == '('
}

func isClosingParentheses(ch byte) bool {
	return ch == ')'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumeric(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphaNumeric(ch byte) bool {
	return isNumeric(ch) || isAlpha(ch)
}

func isMinus(ch byte) bool {
	return ch == '-'
}

func Parse(r *strings.Reader) (Node, error) {
	state := wantOpeningParantheses
	ix := -1
	buf := &bytes.Buffer{}
	nodeStack := make([]Node, 0, 1)
Char:
	for r.Len() > 0 {
		ix++
		ch, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
	State:
		switch state {
		case wantOptionalSpace:
			if isSpace(ch) || isClosingParentheses(ch) {
				if isClosingParentheses(ch) {
					goto Rollup
				}
				continue Char
			}
			state = wantAlphaNumeric
			goto State
		case wantOpeningParantheses:
			if !isOpeningParantheses(ch) {
				return nil, fmt.Errorf("syntax error around ix %d: opening parantheses expected, \"%c\" found", ix, ch)
			}
			state = wantAlpha
			buf.Reset()
			continue Char
		case wantAlpha:
			if isSpace(ch) || isClosingParentheses(ch) {
				lex := buf.String()
				if b, ok := nodeBuilders[lex]; ok {
					nodeStack = append(nodeStack, b())
				} else {
					return nil, fmt.Errorf("unexpected token: %q", lex)
				}
				buf.Reset()
				state = wantOptionalSpace
				if isClosingParentheses(ch) {
					goto Rollup
				}
				continue Char
			}
			if isAlpha(ch) {
				buf.WriteByte(ch)
				continue Char
			}
		case wantNumeric:
			if isSpace(ch) || isClosingParentheses(ch) {
				nodeStack[len(nodeStack)-1].Append(NewIntNode(buf.String()))
				buf.Reset()
				state = wantOptionalSpace
				if isClosingParentheses(ch) {
					goto Rollup
				}
				continue Char
			}
			if isNumeric(ch) {
				buf.WriteByte(ch)
				continue Char
			}
		case wantAlphaNumeric:
			if isSpace(ch) || isClosingParentheses(ch) {
				nodeStack[len(nodeStack)-1].Append(NewVarNode(buf.String()))
				buf.Reset()
				state = wantOptionalSpace
				if isClosingParentheses(ch) {
					goto Rollup
				}
				continue Char
			}
			if (isMinus(ch) || isNumeric(ch)) && buf.Len() == 0 {
				// switch to number-only sequence
				state = wantNumeric
			}
			if (buf.Len() == 0 && isMinus(ch)) || isAlphaNumeric(ch) {
				buf.WriteByte(ch)
				continue Char
			}
			if isOpeningParantheses(ch) {
				buf.Reset()
				state = wantAlpha
				continue Char
			}
		}
		return nil, fmt.Errorf("unexpected token around index %d", ix)
	Rollup:
		if len(nodeStack) > 1 {
			var tail Node
			tail, nodeStack = nodeStack[len(nodeStack)-1], nodeStack[:len(nodeStack)-1]
			nodeStack[len(nodeStack)-1].Append(tail)
			continue Char
		}
	}
	return nodeStack[0], nil
}

func evaluate(expr string) int {
	parsed, err := Parse(strings.NewReader(expr))
	if err != nil {
		panic(fmt.Sprintf("failed to parse expr %q: %s", expr, err))
	}
	res, err := parsed.Eval()
	if err != nil {
		panic(fmt.Sprintf("failed to evaluate expr %q: %s", expr, err))
	}
	return res
}

func main() {
	statements := map[string]int{
		"(add 1 2)":                                  3,
		"(mult 3 (add 2 3))":                         15,
		"(let x 2 (mult x 5))":                       10,
		"(let x 3 x 2 x)":                            2,
		"(let x 1 y 2 x (add x y) (add x y))":        5,
		"(let x 2 (add (let x 3 (let x 4 x)) x))":    6,
		"(let a1 3 b2 (add a1 1) b2)":                4,
		"(let x 2 (mult x (let x 3 y 4 (add x y))))": 14,
		"(let x 7 -12)":                              -12,
	}

	for stmt, wantres := range statements {
		gotres := evaluate(stmt)
		if gotres != wantres {
			panic(fmt.Sprintf("unexpected evaluation result for %q: got: %d, want: %d", stmt, gotres, wantres))
		}
		fmt.Printf("eval %q results to %d\n", stmt, gotres)
	}
}
