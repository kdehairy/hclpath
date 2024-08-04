package parse

import (
	"fmt"
	"log"
	"strconv"

	"github.com/kdehairy/hclpath/v1/lex"
)

type Node string

const (
	Type     Node = "type"
	Attr     Node = "attr"
	Num      Node = "num"
	Str      Node = "str"
	Label    Node = "label"
	Operator Node = "Operator"
)

type Expr interface {
	Print() string
	GetLeft() Expr
	GetRight() Expr
	GetOp() *Op
	GetVal() interface{}
	GetType() Node
}

type Ident struct {
	value string
	ntype Node
}

func (i *Ident) Print() string {
	return i.value
}

func (i *Ident) GetLeft() Expr {
	return nil
}

func (i *Ident) GetRight() Expr {
	return nil
}

func (i *Ident) GetOp() *Op {
	return nil
}

func (i *Ident) GetVal() interface{} {
	return i.value
}

func (i *Ident) GetType() Node {
	return i.ntype
}

type Op string

const (
	SelOp Op = "[]"
	FltOp Op = "{}"
	LblOp Op = ":"
	NstOp Op = "/"
	EqlOp Op = "="
)

func FromToken(tk lex.Token) (op Op) {
	switch tk {
	case lex.NEST:
		op = NstOp
	case lex.SELECT_START:
		op = SelOp
	case lex.FILTER_START:
		op = FltOp
	case lex.NAMED:
		op = LblOp
	case lex.EQUAL:
		op = EqlOp
	}
	return
}

func (o Op) print() string {
	return string(o)
}

type BinOp struct {
	Lhs Expr
	Rhs Expr
	Op  Op
}

func (o *BinOp) Print() string {
	if o.Lhs == nil {
		log.Fatal("Lhs cannot be nil")
	}
	if o.Rhs == nil {
		log.Fatal("Rhs cannot be nil")
	}
	if o.Op == "" {
		log.Fatal("Operator cannot be nil")
	}
	return fmt.Sprintf("(%v-%v-%v)",
		o.Lhs.Print(),
		o.Op.print(),
		o.Rhs.Print())
}

func (o *BinOp) GetLeft() Expr {
	return o.Lhs
}

func (o *BinOp) GetRight() Expr {
	return o.Rhs
}

func (o *BinOp) GetOp() *Op {
	return &o.Op
}

func (o *BinOp) GetVal() interface{} {
	return nil
}

func (o *BinOp) GetType() Node {
	return Operator
}

type NumLt struct {
	value int
}

func (o *NumLt) Print() string {
	return strconv.Itoa(o.value)
}

func (o *NumLt) GetLeft() Expr {
	return nil
}

func (o *NumLt) GetRight() Expr {
	return nil
}

func (o *NumLt) GetOp() *Op {
	return nil
}

func (o *NumLt) GetVal() interface{} {
	return o.value
}

func (o *NumLt) GetType() Node {
	return Num
}

type StrLt struct {
	value string
}

func (o *StrLt) Print() string {
	return o.value
}

func (o *StrLt) GetLeft() Expr {
	return nil
}

func (o *StrLt) GetRight() Expr {
	return nil
}

func (o *StrLt) GetOp() *Op {
	return nil
}

func (o *StrLt) GetVal() interface{} {
	return o.value
}

func (o *StrLt) GetType() Node {
	return Str
}
