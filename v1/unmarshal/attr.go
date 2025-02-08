package unmarshal

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Block struct {
	block *hclsyntax.Block
}

type Attr struct {
	attr *hcl.Attribute
}

func New(block *hclsyntax.Block) Block {
	return Block{
		block: block,
	}
}

func (b Block) GetAttr(name string) (*Attr, error) {
	attrs, _ := b.block.Body.JustAttributes()
	attr, ok := attrs[name]
	if !ok {
		return nil, fmt.Errorf("No attribute with name '%v' found", name)
	}

	return &Attr{attr}, nil
}

func (a *Attr) To(obj interface{}, ctx *hcl.EvalContext) error {
	val, _ := a.attr.Expr.Value(ctx)

	err := gocty.FromCtyValue(val, obj)
	if err != nil {
		return fmt.Errorf("failed to parse value into %v: %v", reflect.TypeOf(obj), err)
	}

	return nil
}
