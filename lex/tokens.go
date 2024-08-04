package lex

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS
	IDENT
	SELECT_START // [
	SELECT_END   // ]
	NAMED        // :
	NEST         // /
	FILTER_START // {
	FILTER_END   // }
	EQUAL        // =
)
