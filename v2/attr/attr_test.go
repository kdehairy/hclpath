package attr

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kdehairy/hclpath/v1"
)

type TestCase struct {
	expected interface{}
	test     func(block *hcl.Block, attr string) (bool, error)
	name     string
	fixture  string
	block    string
	attr     string
}

type AttrObj struct {
	A string `cty:"a"`
	B string `cty:"b"`
	C string `cty:"c"`
}

func TestObjType(t *testing.T) {
	cases := []TestCase{
		{
			name:    "block with object type attribute",
			fixture: "test-1.tf",
			block:   "module",
			attr:    "attr_obj",
			test: func(block *hcl.Block, attr string) (bool, error) {
				attrs := New(block)
				var obj AttrObj
				err := attrs.AsObject(attr, &obj)
				log.Printf("obj: %#v", obj)
				if err != nil {
					return false, fmt.Errorf("failed to parse obj: %v", err)
				}
				if obj.A != "a" {
					return false, nil
				}
				if obj.B != "b" {
					return false, nil
				}
				if obj.C != "c" {
					return false, nil
				}
				return true, nil
			},
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

			if len(blocks) != 1 {
				t.Fatalf("expected 1 but found %v blocks", len(blocks))
			}
			passed, err := tc.test(blocks[0], tc.attr)
			if err != nil {
				t.Fatalf("failed: %v", err)
			}

			if !passed {
				t.Error("Not the expected Object")
			}
		})
	}
}

func TestStringType(t *testing.T) {
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
			block:    "module/block1",
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

			if len(blocks) == 0 {
				t.Fatal("expecting at least 1  block, got 0")
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
