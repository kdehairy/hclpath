package parse

import (
	"fmt"
	"io"

	"github.com/kdehairy/hclpath/lex"
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

func (p *Parser) readSegment() (*Node, error) {
}

func (p *Parser) Parse() (*Node, error) {
	root, err := p.readSegment()
	if err != nil {
		return nil, err
	}

	for {
		tk, lt := p.scanIgnoreWhitespace()
		if tk == lex.EOF {
			return root, nil
		}
		if tk != lex.NEST {
			return nil, fmt.Errorf("Expected '/', but found %v", lt)
		}
		op := NewNode(lt, tk)
		op.Left = root
		root = op

		tk, lt = p.scanIgnoreWhitespace()
		if tk == lex.EOF {
			return nil, fmt.Errorf("Unexpected EOF")
		}
		root.Right = NewNode(lt, tk)
	}
}
