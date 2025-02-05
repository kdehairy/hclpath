package hclpath

import (
	"strings"
	"testing"
)

type TestCase struct {
	name     string
	fixture  string
	test     string
	expected int
}

func TestHclPath(t *testing.T) {
	cases := []TestCase{
		{
			name:     "block without label",
			fixture:  "test-1.tf",
			test:     "terraform",
			expected: 1,
		},
		{
			name:     "attribute value integer",
			fixture:  "test-1.tf",
			test:     "locals{app_version='1'}",
			expected: 1,
		},
		{
			name:     "block with label and label in path",
			fixture:  "test-1.tf",
			test:     "provider:aws",
			expected: 2,
		},
		{
			name:     "child block with label",
			fixture:  "test-1.tf",
			test:     "terraform/backend:s3",
			expected: 1,
		},
		{
			name:     "child block to a block with label",
			fixture:  "test-1.tf",
			test:     "provider:aws/assume_role",
			expected: 1,
		},
		{
			name:     "block with label and attr name",
			fixture:  "test-1.tf",
			test:     "provider:aws{alias}",
			expected: 1,
		},
		{
			name:     "block with label and attr name and value",
			fixture:  "test-1.tf",
			test:     "provider:aws{alias='infra-account'}",
			expected: 1,
		},
		{
			name:     "sub block with label and attr name and value",
			fixture:  "test-1.tf",
			test:     "terraform/backend:s3{region='eu-west-2'}",
			expected: 1,
		},
		{
			name:     "attribute value float",
			fixture:  "test-1.tf",
			test:     "locals{app_float='1.45'}",
			expected: 1,
		},
		{
			name:     "attribute with backslash",
			fixture:  "test-1.tf",
			test:     "provider:aws/assume_role{role_arn='arn:aws:iam::0987654321:role/assumable_role'}",
			expected: 1,
		},
	}

	for _, tc := range cases {
		testName := strings.ReplaceAll(tc.name, " ", "_")
		testName = strings.ToLower(testName)
		t.Run(testName, func(t *testing.T) {
			blocks, err := QueryFile("test_cases/test-1.tf", tc.test)
			if err != nil {
				t.Fatalf("failed to find block: %v", err)
			}

			if len(blocks) != tc.expected {
				t.Errorf("Expected '%v' but found '%v'", tc.expected, len(blocks))
			}
		})
	}
}
