package hclpath

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func readSelector(token string) (string, error) {
	i := strings.IndexRune(token, '[')
	if i == -1 {
		return "", nil
	}
	j := strings.IndexRune(token, ']')
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
	return blockType, label, nil
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
	log.Printf("finding '%v' with label '%v', remaining '%v'...", blockType, blockLabel, rest)
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
