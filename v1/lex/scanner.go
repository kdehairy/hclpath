package lex

import (
	"bufio"
	"bytes"
	"io"
)

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() (ch rune) {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		ch = eof
	}
	return
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) scanWiteSpace() (tk Token, lt string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tk Token, lt string) {
	var buf bytes.Buffer

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLegalIdent(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

func (s *Scanner) scanLiteral() (tk Token, lt string) {
	var buf bytes.Buffer
	startQuote := s.read()
	for {
		if ch := s.read(); ch == startQuote || ch == eof {
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return LITERAL, buf.String()
}

func (s *Scanner) Scan() (tk Token, lt string) {
	ch := s.read()

	if isWhiteSpace(ch) {
		s.unread()
		return s.scanWiteSpace()
	} else if isLegalIdent(ch) {
		s.unread()
		return s.scanIdent()
	}
	switch ch {
	case eof:
		return EOF, ""
	case '/':
		return NEST, string(ch)
	case ':':
		return NAMED, string(ch)
	case '[':
		return SELECT_START, string(ch)
	case ']':
		return SELECT_END, string(ch)
	case '{':
		return FILTER_START, string(ch)
	case '}':
		return FILTER_END, string(ch)
	case '=':
		return EQUAL, string(ch)
	case '\'', '"':
		s.unread()
		return s.scanLiteral()
	}

	return ILLEGAL, string(ch)
}
