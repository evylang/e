package parse

import (
	"fmt"
	"strconv"

	"evylang.dev/e/lex"
)

func ParseExpr(tokens []lex.Token, pratt bool) (n Node, err error) {
	p := parser{
		tokens: tokens,
		scope:  scope{},
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v: token: %v", r, p.cur())
		}
	}()
	n = p.parseExpr(lowest)
	return n, err
}

type precedence int

const (
	lowest precedence = iota
	or
	and
	eq
	lt
	sum
	product
	unary
)

func (p *parser) curPrec() precedence {
	prec, ok := precedences[p.cur().Kind]
	if !ok {
		panic(fmt.Sprintf("no precedence for %v", p.cur()))
	}
	return prec
}

func (p precedence) String() string {
	switch p {
	case lowest:
		return "LOWEST"
	case or:
		return "OR"
	case and:
		return "AND"
	case eq:
		return "EQ"
	case lt:
		return "LT"
	case sum:
		return "SUM"
	case product:
		return "PRODUCT"
	case unary:
		return "UNARY"
	}
	return "🙈"
}

var precedences = map[lex.Kind]precedence{
	lex.Or:  or,
	lex.And: and,
	lex.Eq:  eq, lex.Neq: eq,
	lex.Lt: lt, lex.Gt: lt, lex.Lte: lt, lex.Gte: lt,
	lex.Plus: sum, lex.Minus: sum,
	lex.Slash: product, lex.Star: product,
}

// func (p *parser) parseExpr() Node {
// 	left := p.parseLeft()
// 	if p.atExprEnd() {
// 		return left
// 	}
// 	return p.parseBinaryExpr(left)
// }
//
// 1*1+2
//  Pratt       Naive
//    +            *
//   / \          / \
//  *   2        1   +
// / \              / \
// 1 1             1  2
// 1+2*3
func (p *parser) parseExpr(prec precedence) Node {
	left := p.parseLeft()
	for !p.atExprEnd() && prec < p.curPrec() {
		left = p.parseBinaryExpr(left)
	}
	return left
}

func (p *parser) parseBinaryExpr(left Node) Node {
	op := p.parseBinaryOp()
	prec := precedences[op]
	right := p.parseExpr(prec)
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
	expr := p.parseExpr(lowest)
	p.assertKind(lex.Rparen)
	p.advance() // Advance past )
	return &GroupExpr{expr}
}

func (p *parser) parseUnaryExpr() Node {
	op := p.parseUnaryOp()
	right := p.parseExpr(unary)
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

func (p *parser) atExprEnd() bool {
	kind := p.cur().Kind
	return p.atStmtEnd() || kind == lex.Rparen
}
