package lex

import (
	"fmt"
	"strings"
	"testing"
)

type TestCase struct {
	name     string
	fixture  string
	tk       Token
	expected int
}

func TestMain(t *testing.T) {
	cases := []TestCase{
		{
			name:     "IDENT",
			fixture:  "something",
			tk:       IDENT,
			expected: 1,
		},
		{
			name:     "IDENT with underscore",
			fixture:  "some_thing",
			tk:       IDENT,
			expected: 1,
		},
		{
			name:     "IDENT with dash",
			fixture:  "some-thing",
			tk:       IDENT,
			expected: 1,
		},
		{
			name:     "IDENT with numbers",
			fixture:  "some-thing123",
			tk:       IDENT,
			expected: 1,
		},
		{
			name:     "NEST",
			fixture:  "/",
			tk:       NEST,
			expected: 1,
		},
		{
			name:     "FILTER_START",
			fixture:  "{",
			tk:       FILTER_START,
			expected: 1,
		},
		{
			name:     "FILTER_END",
			fixture:  "}",
			tk:       FILTER_END,
			expected: 1,
		},
		{
			name:     "SELECT_START",
			fixture:  "[",
			tk:       SELECT_START,
			expected: 1,
		},
		{
			name:     "SELECT_END",
			fixture:  "]",
			tk:       SELECT_END,
			expected: 1,
		},
		{
			name:     "NAMED",
			fixture:  ":",
			tk:       NAMED,
			expected: 1,
		},
		{
			name:     "EQUALS",
			fixture:  "=",
			tk:       EQUAL,
			expected: 1,
		},
		{
			name:     "WS",
			fixture:  " ",
			tk:       WS,
			expected: 1,
		},
		{
			name:     "Onw WS for multiple",
			fixture:  "   ",
			tk:       WS,
			expected: 1,
		},
	}
	for _, tc := range cases {
		testName := strings.ReplaceAll(tc.name, " ", "_")
		testName = strings.ToLower(testName)
		t.Run(testName, func(t *testing.T) {
			tokens := scan(tc.fixture)
			if len(tokens) != tc.expected {
				t.Fatalf("Expected %v token found %v", tc.expected, len(tokens))
			}
			err := testForToken(tokens, tc.tk, tc.fixture)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})
	}
}

func scan(str string) (tokens map[string]Token) {
	tokens = make(map[string]Token)
	s := NewScanner(strings.NewReader(str))
	for {
		tk, lt := s.Scan()
		if tk == EOF {
			break
		}
		tokens[lt] = tk
	}
	return
}

func TestSingleIdent(t *testing.T) {
	tokens := scan("test")
	if len(tokens) != 1 {
		t.Fatalf("Expected 1 token found %v", len(tokens))
	}
	tk, ok := tokens["test"]
	if !ok {
		t.Fatalf("Expected a token for string 'test', but found nonw")
	}
	if tk != IDENT {
		t.Fatalf("Expected token of type 'IDENT', but found %v", tk)
	}
}

func testForToken(tokens map[string]Token, expected Token, lt string) error {
	var err error = nil
	tk, ok := tokens[lt]
	if !ok {
		err = fmt.Errorf("Expected a token for string '%v', but found none", lt)
	}
	if tk != expected {
		err = fmt.Errorf("Expected token of type 'IDENT', but found %v", tk)
	}

	return err
}
