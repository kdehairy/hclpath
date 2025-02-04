package hclpath

import (
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kdehairy/hclpath/v2/logging"
)

var logger = logging.NewDefaultLogger()

func Query(b hcl.Body, path string) (hclsyntax.Blocks, error) {
	body := b.(*hclsyntax.Body)
	logger.Debug("Body received", "body", b, "type", reflect.TypeOf(b))
	compilation, err := Compile(path)
	if err != nil {
		return nil, err
	}
	blocks := body.Blocks
	if len(blocks) == 0 {
		logger.Debug("No blocks in the passed body")
		return hclsyntax.Blocks{}, nil
	}
	return compilation.Exec(blocks)
}
