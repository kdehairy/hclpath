package parse

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kdehairy/hclpath/v2/lex"
)

type Parser struct {
	s   *lex.Scanner
	buf struct {
		lt          string
		tk          lex.Token
		isUnscanned bool
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: lex.NewScanner(r),
	}
}

func (p *Parser) scan() (tk lex.Token, lt string) {
	if p.buf.isUnscanned {
		p.buf.isUnscanned = false
		return p.buf.tk, p.buf.lt
	}

	tk, lt = p.s.Scan()

	p.buf.tk, p.buf.lt = tk, lt

	return
}

func (p *Parser) unscan() {
	p.buf.isUnscanned = true
}

func (p *Parser) scanIgnoreWhitespace() (tk lex.Token, lt string) {
	tk, lt = p.scan()
	if tk == lex.WS {
		tk, lt = p.scan()
	}
	return
}

func (p *Parser) peek() (tk lex.Token) {
	tk, _ = p.scanIgnoreWhitespace()
	p.unscan()
	return
}

func (p *Parser) expect(expected lex.Token) (res bool) {
	res = p.peek() == expected
	return
}

func (p *Parser) consume(expected lex.Token) (bool, error) {
	tk, _ := p.scanIgnoreWhitespace()
	if tk != expected {
		return false, fmt.Errorf("expected '%v' but found '%v'", expected, tk)
	}
	return true, nil
}

func (p *Parser) parseIdent() (Expr, error) {
	tk, lt := p.scanIgnoreWhitespace()
	if tk != lex.IDENT {
		return nil, fmt.Errorf("expected identifier found %v", tk)
	}
	return &Ident{value: lt}, nil
}

func (p *Parser) parseFilter() (Expr, error) {
	lhs, err := p.parseIdent()
	lhs.(*Ident).ntype = Attr
	if err != nil {
		return nil, fmt.Errorf("failed to parser filter: %v", err)
	}
	if ok := p.expect(lex.EQUAL); ok {
		p.consume(lex.EQUAL)
		rhs, err := p.parseLiteral()
		if err != nil {
			return nil, fmt.Errorf("failed to parser filter: %v", err)
		}
		return &BinOp{
			Lhs: lhs,
			Rhs: rhs,
			Op:  EqlOp,
		}, nil
	} else {
		return lhs, nil
	}
}

func (p *Parser) parseLiteral() (Expr, error) {
	tk, lt := p.scanIgnoreWhitespace()
	if tk != lex.LITERAL {
		return nil, fmt.Errorf("expected %v found %v", lex.LITERAL, tk)
	}

	return &StrLt{
		value: lt,
	}, nil
}

func (p *Parser) parseNum() (Expr, error) {
	tk, lt := p.scanIgnoreWhitespace()
	if tk != lex.IDENT {
		return nil, fmt.Errorf("expected integer found %v", tk)
	}
	i, err := strconv.Atoi(lt)
	if err != nil {
		return nil, fmt.Errorf("expected integer, but found '%v'", lt)
	}

	return &NumLt{
		value: i,
	}, nil
}

func (p *Parser) Parse() (Expr, error) {
	lhs, err := p.parseIdent()
	lhs.(*Ident).ntype = Type
	if err != nil {
		return nil, fmt.Errorf("syntax error: %v", err)
	}

	for p.peek().IsOperator() {
		tk, _ := p.scanIgnoreWhitespace()

		var rhs Expr
		var err error
		op := FromToken(tk)
		switch tk {
		case lex.NEST:
			rhs, err = p.parseIdent()
			rhs.(*Ident).ntype = Type
		case lex.NAMED:
			rhs, err = p.parseIdent()
			rhs.(*Ident).ntype = Label
		case lex.FILTER_START:
			rhs, err = p.parseFilter()
			if ok, error := p.consume(lex.FILTER_END); !ok {
				return nil, error
			}
		case lex.SELECT_START:
			rhs, err = p.parseNum()
			if ok, error := p.consume(lex.SELECT_END); !ok {
				return nil, error
			}
		}
		if err != nil {
			return nil, fmt.Errorf("syntax error: %v", err)
		}
		lhs = &BinOp{
			Op:  op,
			Lhs: lhs,
			Rhs: rhs,
		}
	}

	return lhs, nil
}
