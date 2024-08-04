package parse

import "github.com/kdehairy/hclpath/lex"

type Node struct {
	Left    *Node
	Right   *Node
	Literal string
	Tk      lex.Token
}

func NewNode(literal string, tk lex.Token) *Node {
	return &Node{Left: nil, Right: nil, Literal: literal, Tk: tk}
}
