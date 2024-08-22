package hclpath

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/kdehairy/hclpath/v2/cmpval"
	"github.com/kdehairy/hclpath/v2/logging"
	"github.com/kdehairy/hclpath/v2/parse"
)

type (
	execFunc func(hclsyntax.Blocks) (hclsyntax.Blocks, error)
	evalFunc func(hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error)
)

type Compilation struct {
	Exec execFunc
}

type evaluation struct {
	Do evalFunc
}

var logger = logging.NewDefaultLogger(slog.LevelDebug)

func Compile(path string) (*Compilation, error) {
	logger.Info("Recieved path", "path", path)
	p := parse.NewParser(strings.NewReader(path))
	expr, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("syntax error: %v", err)
	}
	logger.Debug("AST", "expr", expr.Print())

	eval, _, err := evaluate(expr)
	if err != nil {
		return nil, err
	}

	Compilation := &Compilation{
		Exec: func(b hclsyntax.Blocks) (blocks hclsyntax.Blocks, err error) {
			logger.Debug("Executing Compilation...")
			blocks, _, err = eval.Do(b)
			return
		},
	}

	logger.Debug("Compilation Complete", "Compilation", Compilation)

	return Compilation, nil
}

func evaluate(expr parse.Expr) (*evaluation, interface{}, error) {
	logger.Debug(">> Evaluating Expr:", "AST", expr.Print(), "type", expr.GetType())
	var lhs *evaluation
	var rhs *evaluation
	var lvalue interface{}
	var rvalue interface{}
	var value interface{}
	var self *evaluation
	var err error
	logger.Debug(">> Start evaluating sides")
	if expr.GetLeft() != nil {
		lhs, lvalue, err = evaluate(expr.GetLeft())
		if err != nil {
			return nil, nil, err
		}
	}

	if expr.GetRight() != nil {
		rhs, rvalue, err = evaluate(expr.GetRight())
		if err != nil {
			return nil, nil, err
		}
	}
	logger.Debug(">> End evaluating sides")

	self = &evaluation{}
	if expr.GetOp() != nil {
		logger.Debug("Expr operator", "Op", *expr.GetOp())
		switch *expr.GetOp() {
		case parse.NstOp:
			self.Do = func(b hclsyntax.Blocks) (blocks hclsyntax.Blocks, value interface{}, err error) {
				blocks, _, err = lhs.Do(b)
				if err != nil {
					return nil, nil, err
				}
				if blocks == nil {
					return nil, nil, errors.New("unexpected null blocks on left of operator")
				}
				candidates := hclsyntax.Blocks{}
				for _, b := range blocks {
					if res, _, e := rhs.Do(b.Body.Blocks); len(res) > 0 && e == nil {
						candidates = append(candidates, res...)
					} else if e != nil {
						return nil, nil, e
					}
				}
				return candidates, nil, nil
			}
		case parse.FltOp, parse.LblOp:
			self.Do = func(b hclsyntax.Blocks) (blocks hclsyntax.Blocks, value interface{}, err error) {
				blocks, _, err = lhs.Do(b)
				if err != nil {
					return nil, nil, err
				}
				if blocks == nil {
					return nil, nil, errors.New("unexpected null blocks on left of operator")
				}
				return rhs.Do(blocks)
			}
		case parse.EqlOp:
			self.Do = func(b hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error) {
				blocks, _, err := lhs.Do(b)
				if err != nil {
					return nil, nil, err
				}
				if lvalue == nil {
					return nil, nil, errors.New("expected lvalue, but found none")
				}
				if rvalue == nil {
					return nil, nil, errors.New("expected rvalue, but found none")
				}

				attrName, ok := lvalue.(string)
				if !ok {
					return nil, nil, fmt.Errorf("expected string lvalue, but found '%v'", lvalue)
				}
				attrValue, ok := rvalue.(string)
				if !ok {
					return nil, nil, fmt.Errorf("expected string rvalue, but found '%v'", rvalue)
				}

				return filter(blocks, attrName, attrValue)
			}
		case parse.SelOp:
			self.Do = func(b hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error) {
				// TODO: find blocks
				return nil, nil, errors.ErrUnsupported
			}
		}
	} else {
		logger.Debug("No operator")
		switch expr.GetType() {
		case parse.Type:
			value = expr.GetVal()
			self.Do = func(b hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error) {
				logger.Debug("Evaluating 'type' Node", "expr", expr.Print())
				val := expr.GetVal()
				if val == nil {
					return nil, nil, errors.New("expected block type, but found none")
				}
				name, ok := val.(string)
				if !ok {
					return nil, nil, fmt.Errorf("expected string, but found '%v'", name)
				}
				blocks := findBlocksByType(b, name)
				return blocks, val, nil
			}
		case parse.Label:
			value = expr.GetVal()
			self.Do = func(b hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error) {
				logger.Debug("Evaluating 'label' Node", "expr", expr.Print())
				val := expr.GetVal()
				if val == nil {
					return nil, nil, errors.New("expected block label, but found none")
				}
				name, ok := val.(string)
				if !ok {
					return nil, nil, fmt.Errorf("expected string, but found '%v'", name)
				}
				blocks := findBlocksByLabel(b, name)
				return blocks, val, nil
			}
		case parse.Attr:
			value = expr.GetVal()
			self.Do = func(b hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error) {
				logger.Debug("Evaluating 'attr' Node", "expr", expr.Print())
				val := expr.GetVal()
				if val == nil {
					return nil, nil, errors.New("expected block label, but found none")
				}
				name, ok := val.(string)
				if !ok {
					return nil, nil, fmt.Errorf("expected string, but found '%v'", name)
				}
				blocks, err := findBlocksByAttr(b, name)
				return blocks, val, err
			}
		case parse.Num, parse.Str:
			value = expr.GetVal()
			self.Do = func(b hclsyntax.Blocks) (hclsyntax.Blocks, interface{}, error) {
				logger.Debug("Evaluating 'literal' Node", "expr", expr.Print())
				return nil, expr.GetVal(), nil
			}
		}
	}

	logger.Debug("Finished Evaluation", "expr", expr.Print())
	return self, value, nil
}

func findBlocksByAttr(blocks hclsyntax.Blocks, name string) (hclsyntax.Blocks, error) {
	var candidates hclsyntax.Blocks = []*hclsyntax.Block{}
	for _, b := range blocks {
		attrs, _ := b.Body.JustAttributes()
		if attrs == nil {
			return nil, fmt.Errorf("failed to read attributes from block '%v'", b.Type)
		}
		if len(attrs) == 0 {
			continue
		}
		for _, a := range attrs {
			if a.Name == name {
				candidates = append(candidates, b)
				break
			}
		}
	}
	return candidates, nil
}

func findBlocksByLabel(blocks hclsyntax.Blocks, name string) hclsyntax.Blocks {
	logger.Debug("### findBlocksByLabel")
	var candidates hclsyntax.Blocks = []*hclsyntax.Block{}
	logger.Debug("### Blocks", "count", len(blocks))
	for _, b := range blocks {
		logger.Debug("### Labels", "block", b.Type, "count", len(b.Labels))
		for _, l := range b.Labels {
			if l == name {
				logger.Debug("Found block with label", "block", b.Type, "label", l)
				candidates = append(candidates, b)
				break
			}
		}
	}
	return candidates
}

func findBlocksByType(blocks hclsyntax.Blocks, name string) hclsyntax.Blocks {
	logger.Info("Finding Block by type...", "type", name)
	var candidates hclsyntax.Blocks = []*hclsyntax.Block{}
	logger.Debug("### Blocks", "count", len(blocks))
	for _, b := range blocks {
		logger.Debug("Examining block", "block", b.Type)
		if b.Type == name {
			logger.Debug("Found block", "block", b.Type, "label", b.Type)
			candidates = append(candidates, b)
		}
	}
	return candidates
}

func filter(blocks hclsyntax.Blocks, attrName string, attrValue string) (hclsyntax.Blocks, interface{}, error) {
	var candidateBlocks hclsyntax.Blocks

	for _, b := range blocks {
		attrs, _ := b.Body.JustAttributes()
		if attrs == nil {
			return nil, nil, errors.New("failed to read attributes")
		}
		if len(attrs) == 0 {
			continue
		}

		for _, a := range attrs {
			if a.Name != attrName {
				continue
			}
			if attrValue == "" {
				candidateBlocks = append(candidateBlocks, b)
				continue
			}

			val, _ := a.Expr.Value(nil)
			equals, err := cmpval.IsEqual(val, attrValue)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to test equality: %v", err)
			}
			if equals {
				candidateBlocks = append(candidateBlocks, b)
			}
		}
	}
	return candidateBlocks, nil, nil
}
