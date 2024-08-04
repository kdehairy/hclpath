package parse

import (
	"strings"
	"testing"
)

type TestCase struct {
	name     string
	fixture  string
	expected string
}

func TestParser(t *testing.T) {
	cases := []TestCase{
		{
			name:     "first/second",
			fixture:  "first/second",
			expected: "(first-/-second)",
		},
		{
			name:     "first:label",
			fixture:  "first:label",
			expected: "(first-:-label)",
		},
		{
			name:     "first:label/second",
			fixture:  "first:label/second",
			expected: "((first-:-label)-/-second)",
		},
		{
			name:     "first/second:label",
			fixture:  "first/second:label",
			expected: "((first-/-second)-:-label)",
		},
		{
			name:     "first{attr=val}",
			fixture:  "first{attr='val'}",
			expected: "(first-{}-(attr-=-val))",
		},
		{
			name:     "first{attr=val}/second",
			fixture:  "first{attr='val'}/second",
			expected: "((first-{}-(attr-=-val))-/-second)",
		},
		{
			name:     "first[12]",
			fixture:  "first[12]",
			expected: "(first-[]-12)",
		},
		{
			name:     "first[12]/second",
			fixture:  "first[12]/second",
			expected: "((first-[]-12)-/-second)",
		},
		{
			name:     "first/second[12]",
			fixture:  "first/second[12]",
			expected: "((first-/-second)-[]-12)",
		},
		{
			name:     "first:label{attr}",
			fixture:  "first:label{attr}",
			expected: "((first-:-label)-{}-attr)",
		},
	}

	for _, tc := range cases {
		testName := strings.ReplaceAll(tc.name, " ", "_")
		testName = strings.ToLower(testName)
		t.Run(testName, func(t *testing.T) {
			p := NewParser(strings.NewReader(tc.fixture))
			expr, err := p.Parse()
			if err != nil {
				t.Fatalf("%v", err)
			}
			found := expr.Print()
			if tc.expected != found {
				t.Fatalf("expected '%v' but found '%v'", tc.expected, found)
			}
		})
	}
}
