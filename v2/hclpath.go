package hclpath

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kdehairy/hclpath/v2/logging"
)

var logger = logging.NewDefaultLogger()

func QueryFile(file string, path string) (hclsyntax.Blocks, error) {
	hclParser := hclparse.NewParser()
	hclFile, _ := hclParser.ParseHCLFile("test_cases/test-1.tf")
	if hclFile == nil {
		return nil, fmt.Errorf("failed to parse file '%v'", file)
	}
	return Query(hclFile.Body, path)
}

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
