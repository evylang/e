package parse

import (
	"fmt"

	"evylang.dev/e/lex"
)

func ParseProg(tokens []lex.Token) (prog *Prog, err error) {
	p := parser{
		tokens: tokens,
		scope:  scope{},
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v: token: %v", r, p.cur())
		}
	}()
	prog = p.parseProg()
	return prog, err
}

type parser struct {
	tokens []lex.Token
	pos    int
	scope  scope
}

type scope map[string]*Var

func (s scope) add(v *Var) {
	if _, ok := s[v.Name]; ok {
		panic("redeclaration " + v.Name)
	}
	s[v.Name] = v
}

func (s scope) get(name string) *Var {
	v, ok := s[name]
	if !ok {
		panic("undeclared " + name)
	}
	return v
}

func (p *parser) parseProg() *Prog {
	var stmts []Node
	p.eatWS()
	for p.cur().Kind != lex.Eof {
		stmt := p.parseStmt()
		stmts = append(stmts, stmt)
		p.eatWS()
	}
	return &Prog{Stmts: stmts}
}

func (p *parser) eatWS() {
	k := p.cur().Kind
	for k == lex.Nl || k == lex.Comment {
		p.advance()
		k = p.cur().Kind
	}
}

func (p *parser) parseStmt() Node {
	switch p.cur().Kind {
	case lex.If:
		return p.parseIfStmt()
	case lex.Ident:
		switch p.peek().Kind {
		case lex.Assign:
			return p.parseAssignStmt()
		case lex.Colon:
			return p.parseDeclStmt()
		}
	}
	panic("bad statement")
}

func (p *parser) parseIfStmt() Node {
	return nil
}

func (p *parser) parseAssignStmt() Node {
	// ident "=" expr .
	varName := p.cur().Lit
	target := p.scope.get(varName)
	p.advance() // advance past IDENT
	p.advance() // advance past ASSIGN
	val := p.parseExpr()
	assertTypes(target, lex.Assign, val)
	return &Assign{target, val}
}

func (p *parser) parseDeclStmt() Node {
	varName := p.parseIdent()
	p.assertKind(lex.Colon)
	p.advance() // advance past COLON
	t := p.parseType()
	p.assertStmtEnd()
	v := &Var{Name: varName, T: t}
	p.scope.add(v)
	return &Decl{Var: v}
}

func (p *parser) parseType() Type {
	ident := p.parseIdent()
	p.assertStmtEnd()
	switch ident {
	case "num":
		return NumType
	case "string":
		return StringType
	case "bool":
		return BoolType
	}
	panic("unknown type " + ident)
}

func (p *parser) parseIdent() string {
	if p.cur().Kind != lex.Ident {
		panic("cannot parse ident")
	}
	ident := p.cur().Lit
	p.advance() // advance past IDENT
	return ident
}

func (p *parser) advance() {
	p.pos++
}

func (p *parser) cur() lex.Token {
	if p.pos >= len(p.tokens) {
		return lex.Token{Kind: lex.Eof}
	}
	return p.tokens[p.pos]
}

func (p *parser) peek() lex.Token {
	if p.pos+1 >= len(p.tokens) {
		return lex.Token{Kind: lex.Eof}
	}
	return p.tokens[p.pos+1]
}

func (p *parser) assertStmtEnd() {
	if !p.atStmtEnd() {
		panic("expected end of statement")
	}
}

func (p *parser) atStmtEnd() bool {
	kind := p.cur().Kind
	return kind == lex.Nl || kind == lex.Eof || kind == lex.Comment
}

func (p *parser) assertKind(w lex.Kind) {
	g := p.cur().Kind
	if w != g {
		format := "token: want: %v, got: %v"
		panic(fmt.Sprintf(format, w, g))
	}
}

func assertTypes(left Node, op lex.Kind, right Node) {
	leftT := left.Type()
	rightT := right.Type()
	if leftT != rightT {
		panic(fmt.Sprintf("op: %v. types not equal: %v != %v", op, leftT, rightT))
	}
	switch op {
	case lex.Minus, lex.Star, lex.Slash:
		if leftT != NumType {
			panic(fmt.Sprintf("want: num, got: %s for %s operator", leftT, op))
		}
	case lex.And, lex.Or:
		if leftT != BoolType {
			panic(fmt.Sprintf("want: bool, got: %s for %s operator", leftT, op))
		}
	case lex.Plus, lex.Lt, lex.Gt, lex.Lte, lex.Gte:
		if leftT != NumType && rightT != StringType {
			panic(fmt.Sprintf("want: num or string, got: %s for %s operator", leftT, op))
		}
	}
	// lex.Eq and lex.Neq work for all types
}
