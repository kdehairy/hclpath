package lex

type Token string

const (
	ILLEGAL      Token = "Illegal"
	EOF          Token = "EOF"
	WS           Token = "WhiteSpace"
	IDENT        Token = "Identity"
	SELECT_START Token = "["
	SELECT_END   Token = "]"
	NAMED        Token = ":"
	NEST         Token = "/"
	FILTER_START Token = "{"
	FILTER_END   Token = "}"
	EQUAL        Token = "="
	QUOTE        Token = "'"
	DQUOTE       Token = "\""
	LITERAL      Token = "literal"
)

func (t Token) IsOperator() bool {
	return t == SELECT_START ||
		t == NAMED ||
		t == NEST ||
		t == FILTER_START ||
		t == EQUAL
}
