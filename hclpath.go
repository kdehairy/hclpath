package hclpath

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func readSelector(token string) (string, error) {
	i := strings.Index(token, "[")
	if i == -1 {
		return "", nil
	}
	j := strings.LastIndex(token, "]")
	if j == -1 {
		return "", errors.New("no matching ']' found")
	}

	return token[i+1 : j], nil
}

func readBlockHeader(token string) (blockType string, blockLabels string, err error) {
	blockType, rest, _ := strings.Cut(token, ":")
	if blockType == "" {
		return "", "", errors.New("no type token found")
	}

	label, _, _ := strings.Cut(rest, "[")
	blockType, _, _ = strings.Cut(blockType, "[")
	return blockType, label, nil
}

func filterBlocks(blocks hcl.Blocks, selector string) (hcl.Blocks, error) {
	attrName, attrValue, _ := strings.Cut(selector, "=")
	log.Printf("attrName: %v, attrValue: %v", attrName, attrValue)

	var candidateBlocks hcl.Blocks

	for _, b := range blocks {
		attrs, _ := b.Body.JustAttributes()
		if attrs == nil {
			return nil, errors.New("failed to read attributes")
		}
		if len(attrs) == 0 {
			continue
		}

		for _, a := range attrs {
			if a.Name != attrName {
				continue
			}
			if attrValue == "" {
				candidateBlocks = append(candidateBlocks, b)
				continue
			}

			val, _ := a.Expr.Value(nil)
			equals, err := isEqual(val, attrValue)
			if err != nil {
				return nil, fmt.Errorf("failed to test equality: %v", err)
			}
			if equals {
				candidateBlocks = append(candidateBlocks, b)
			}
		}
	}

	return candidateBlocks, nil
}

func FindBlocks(b hcl.Body, path string) (hcl.Blocks, error) {
	token, rest, _ := strings.Cut(path, "/")

	blockType, blockLabel, err := readBlockHeader(token)
	if err != nil {
		return nil, fmt.Errorf("invalid path '%v':%v", path, err)
	}

	var labels []string
	if blockLabel != "" {
		labels = []string{blockLabel}
	}
	log.Printf("finding '%v' with labels '%v', remaining '%v'...", blockType, labels, rest)
	schema := hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       blockType,
				LabelNames: labels,
			},
		},
	}
	content, _, _ := b.PartialContent(&schema)
	if content == nil {
		log.Fatal("Failed to read content of Hcl file")
	}
	blocks := content.Blocks
	if len(blocks) == 0 {
		log.Printf("No blocks found")
		return hcl.Blocks{}, nil
	}

	selector, err := readSelector(token)
	if err != nil {
		return nil, fmt.Errorf("failed to read selectors: %v", err)
	}

	if selector != "" {
		filteredBlocks, err := filterBlocks(blocks, selector)
		if err != nil {
			return nil, fmt.Errorf("failed to filter by selector: %v", err)
		}
		blocks = filteredBlocks
	}

	var candidateBlocks hcl.Blocks
	if rest == "" {
		candidateBlocks = blocks
		log.Printf("found %v candidates", len(candidateBlocks))
		return candidateBlocks, nil
	}

	for _, block := range blocks {
		newFound, err := FindBlocks(block.Body, rest)
		if err != nil {
			return nil, fmt.Errorf("invalid path '%v':%v", path, err)
		}
		candidateBlocks = append(candidateBlocks, newFound[:]...)
	}
	return candidateBlocks, nil
}
