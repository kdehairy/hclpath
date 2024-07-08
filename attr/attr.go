package attr

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Attrs struct {
	block *hcl.Block
}

func New(block *hcl.Block) Attrs {
	return Attrs{
		block: block,
	}
}

func getAttr(body hcl.Body, name string) (hcl.Attribute, error) {
	schema := hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name: name,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{},
	}
	content, _, _ := body.PartialContent(&schema)
	if content == nil {
		return hcl.Attribute{}, errors.New("failed to read body")
	}
	if len(content.Attributes) != 1 {
		return hcl.Attribute{}, fmt.Errorf("found %v but expected 1 attribute with name %v",
			len(content.Attributes), name)
	}

	return *content.Attributes[name], nil
}

func (a Attrs) AsString(name string) (string, error) {
	attr, err := getAttr(a.block.Body, name)
	if err != nil {
		return "", fmt.Errorf("failed to get attribute '%v': %v", name, err)
	}
	val, _ := attr.Expr.Value(nil)
	if val.Type() != cty.String {
		return "", fmt.Errorf("attribute '%v' is of type '%v'", name, val.Type().FriendlyName())
	}

	return val.AsString(), nil
}

func (a Attrs) AsObject(name string, obj interface{}) error {
	attr, err := getAttr(a.block.Body, name)
	if err != nil {
		return fmt.Errorf("failed to get attribute '%v': %v", name, err)
	}
	val, _ := attr.Expr.Value(nil)
	if !val.Type().IsObjectType() {
		return fmt.Errorf("attribute '%v' is of type '%v'", name, val.Type().FriendlyName())
	}

	err = gocty.FromCtyValue(val, obj)
	if err != nil {
		return fmt.Errorf("failed to parse value into %v", reflect.TypeOf(obj))
	}

	return nil
}
