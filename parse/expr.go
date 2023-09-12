package parse

import (
	"fmt"
	"strconv"

	"evylang.dev/e/lex"
)

func (p *parser) parseExpr() Node {
	left := p.parseLeft()
	if p.atExprEnd() {
		return left
	}
	return p.parseBinaryExpr(left)
}

func (p *parser) parseBinaryExpr(
	left Node) Node {
	op := p.parseBinaryOp()
	right := p.parseExpr()
	assertTypes(left, op, right)
	return &BinaryExpr{left, op, right}
}

func (p *parser) parseLeft() Node {
	switch {
	case p.cur().IsUnaryOp():
		return p.parseUnaryExpr()
	case p.cur().Kind == lex.Lparen:
		return p.parseGroupExpr()
	}
	return p.parseOperand()
}

func (p *parser) parseGroupExpr() Node {
	p.assertKind(lex.Lparen)
	p.advance() // Advance past (
	expr := p.parseExpr()
	p.assertKind(lex.Rparen)
	p.advance() // Advance past )
	return &GroupExpr{expr}
}

func (p *parser) parseUnaryExpr() Node {
	op := p.parseUnaryOp()
	right := p.parseExpr()
	return &UnaryExpr{Op: op, Right: right}
}

func (p *parser) parseBinaryOp() lex.Kind {
	if !p.cur().IsBinaryOp() {
		panic(fmt.Sprintf("invalid binary operator %v", p.cur()))
	}
	k := p.cur().Kind
	p.advance() // Advance past OP
	return k
}

func (p *parser) parseUnaryOp() lex.Kind {
	if !p.cur().IsUnaryOp() {
		panic(fmt.Sprintf("invalid unary operator %v", p.cur()))
	}
	k := p.cur().Kind
	p.advance() // Advance past OP
	return k
}

func (p *parser) parseOperand() Node {
	cur := p.cur()
	p.advance()
	switch cur.Kind {
	case lex.StringLit:
		return &StringLit{Val: cur.Lit}
	case lex.NumLit:
		f := parseNum(cur.Lit)
		return &NumLit{Val: f}
	case lex.Ident:
		return p.scope.get(cur.Lit)
	case lex.True:
		return &BoolLit{Val: true}
	case lex.False:
		return &BoolLit{Val: false}
	}
	panic("💥")
}

func parseNum(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic("invalid num: " + s)
	}
	return f
}

// func (p *parser) parseOperand() Node {
// 	cur := p.cur()
// 	p.advance()
// 	switch cur.Kind {
// 	case lex.StringLit:
// 		return &StringLit{Val: cur.Lit}
// 	case lex.NumLit: // ...
// 	}
// 	panic("💥")
// }

// case lex.NumLit:
// 	f := parseFloat(cur.Lit)
// 	return &NumLit{Val: f}
// case lex.True:
// 	return &BoolLit{Val: true}
// case lex.False:
// 	return &BoolLit{Val: false}
// case lex.Ident:
// 	return p.scope.get(cur.Lit)

func (p *parser) atExprEnd() bool {
	kind := p.cur().Kind
	return p.atStmtEnd() || kind == lex.Rparen
}
