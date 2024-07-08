package attr

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kdehairy/hclpath"
)

type TestCase struct {
	expected interface{}
	name     string
	fixture  string
	block    string
	attr     string
}

func TestAttr(t *testing.T) {
	cases := []TestCase{
		{
			name:     "block with both sub-blocks and attrs",
			fixture:  "test-1.tf",
			block:    "module",
			attr:     "attr01",
			expected: "x",
		},
		{
			name:     "sub-block",
			fixture:  "test-1.tf",
			block:    "module:test/block1",
			attr:     "attr11",
			expected: "a",
		},
	}

	for _, tc := range cases {
		testName := strings.ReplaceAll(tc.name, " ", "_")
		testName = strings.ToLower(testName)
		t.Run(testName, func(t *testing.T) {
			hclParser := hclparse.NewParser()
			hclFile, _ := hclParser.ParseHCLFile("test_cases/test-1.tf")
			if hclFile == nil {
				t.Fatalf("failed to parse hcl file")
			}
			body := hclFile.Body

			blocks, err := hclpath.FindBlocks(body, tc.block)
			if err != nil {
				t.Fatalf("failed to find block:%v", err)
			}

			for _, b := range blocks {
				attrs := New(b)
				str, err := attrs.AsString(tc.attr)
				if err != nil {
					t.Fatalf("error while reading attribute: %v", err)
				}
				if str != tc.expected.(string) {
					t.Errorf("expected '%v' but found '%v'", tc.expected.(string), str)
				}
			}
		})
	}
}
