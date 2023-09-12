package lex

import (
	"unicode"
)

func Tokenize(input string) []Token {
	l := lexer{input: []rune(input)}
	var tokens []Token
	for token := l.next(); token.Kind != Eof; token = l.next() {
		tokens = append(tokens, token)
	}
	return tokens
}

type Kind int

const (
	Eof Kind = iota
	Colon
	Assign
	Lparen
	Rparen
	Nl

	// Operators
	Plus
	Minus
	Star
	Slash
	Eq
	Neq
	Lt
	Gt
	Lte
	Gte
	Bang

	// Literal
	StringLit
	NumLit
	Comment
	Ident

	// Keywords
	If
	While
	End
	True
	False
	And
	Or
)

func (k Kind) String() string {
	return kindStrings[k]
}

func (k Kind) Repr() string {
	return reprStrings[k]
}

func (k Kind) IsComparisonOp() bool {
	return k == Eq || k == Neq || k == Lt || k == Gt || k == Lte || k == Gte
}

type Token struct {
	Kind Kind
	Lit  string
}

func (t Token) IsUnaryOp() bool {
	return t.Kind == Bang || t.Kind == Minus
}

func (t Token) IsBinaryOp() bool {
	k := t.Kind
	return k == Plus || k == Minus || k == Star || k == Slash || k == And || k == Or || k.IsComparisonOp()
}

type lexer struct {
	input []rune
	pos   int
}

func (l *lexer) cur() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *lexer) advance() {
	l.pos++
}

func (l *lexer) next() Token {
	l.eatWhiteSpace()
	cur := l.cur()
	l.advance()
	peek := l.cur()
	switch cur {
	case 0:
		return Token{Kind: Eof}
	case '\n':
		return Token{Kind: Nl}
	case '+':
		return Token{Kind: Plus}
	case '-':
		return Token{Kind: Minus}
	case '*':
		return Token{Kind: Star}
	case ':':
		return Token{Kind: Colon}
	case '(':
		return Token{Kind: Lparen}
	case ')':
		return Token{Kind: Rparen}
	case '=':
		if peek == '=' {
			l.advance()
			return Token{Kind: Eq}
		}
		return Token{Kind: Assign}
	case '>':
		if peek == '=' {
			l.advance()
			return Token{Kind: Gte}
		}
		return Token{Kind: Gt}
	case '<':
		if peek == '=' {
			l.advance()
			return Token{Kind: Lte}
		}
		return Token{Kind: Lt}
	case '/':
		if peek == '/' {
			l.advance()
			return l.comment()
		}
		return Token{Kind: Slash}
	case '!':
		if peek == '=' {
			l.advance()
			return Token{Kind: Neq}
		}
		return Token{Kind: Bang}
	case '"':
		return l.stringLit()
	}
	if unicode.IsDigit(cur) {
		return l.numLit(cur)
	}
	if isLetter(cur) {
		s := l.ident(cur)
		if kind, ok := keywords[s]; ok {
			return Token{Kind: kind}
		}
		return Token{Kind: Ident, Lit: s}
	}
	panic("unknown char" + string(cur))
}

func (l *lexer) eatWhiteSpace() {
	for l.cur() == ' ' || l.cur() == '\t' {
		l.advance()
	}
}

func (l *lexer) comment() Token {
	var runes []rune
	for l.cur() != '\n' && l.cur() != 0 {
		runes = append(runes, l.cur())
		l.advance()
	}
	lit := string(runes)
	return Token{Kind: Comment, Lit: lit}
}

func (l *lexer) numLit(r rune) Token {
	runes := []rune{r}
	for unicode.IsDigit(l.cur()) {
		runes = append(runes, l.cur())
		l.advance()
	}
	lit := string(runes)
	return Token{Kind: NumLit, Lit: lit}
}

func (l *lexer) stringLit() Token {
	runes := []rune{}
	for l.cur() != '\n' && l.cur() != 0 && l.cur() != '"' {
		runes = append(runes, l.cur())
		l.advance()
	}
	if l.cur() != '"' {
		panic("unterminated string")
	}
	l.advance()
	lit := string(runes)
	return Token{Kind: StringLit, Lit: lit}
}

func (l *lexer) ident(r rune) string {
	runes := []rune{r}
	for unicode.IsDigit(l.cur()) || isLetter(l.cur()) {
		runes = append(runes, l.cur())
		l.advance()
	}
	return string(runes)
}

func isLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

var kindStrings = map[Kind]string{
	Eof:       "EOF",
	Ident:     "IDENT",
	Colon:     "COLON",
	Assign:    "ASSIGN",
	Lparen:    "LPAREN",
	Rparen:    "RPAREN",
	Nl:        "NL",
	Plus:      "PLUS",
	Minus:     "MINUS",
	Star:      "STAR",
	Slash:     "SLASH",
	Eq:        "EQ",
	Neq:       "NEQ",
	Lt:        "LT",
	Gt:        "GT",
	Lte:       "LTE",
	Gte:       "GTE",
	Bang:      "BANG",
	StringLit: "STRIN_LIT",
	NumLit:    "NUM_LIT",
	Comment:   "COMMENT",
	If:        "IF",
	While:     "WHILE",
	End:       "END",
	True:      "TRUE",
	False:     "FALSE",
	And:       "AND",
	Or:        "OR",
}

var reprStrings = map[Kind]string{
	Eof:       "EOF",
	Ident:     "IDENT",
	Colon:     "COLON",
	Assign:    "ASSIGN",
	Lparen:    "LPAREN",
	Rparen:    "RPAREN",
	Nl:        "NL",
	Plus:      "+",
	Minus:     "-",
	Star:      "*",
	Slash:     "/",
	Eq:        "==",
	Neq:       "!=",
	Lt:        "<",
	Gt:        ">",
	Lte:       "<=",
	Gte:       ">=",
	Bang:      "!",
	StringLit: "STRIN_LIT",
	NumLit:    "NUM_LIT",
	Comment:   "COMMENT",
	If:        "if",
	While:     "while",
	End:       "end",
	True:      "true",
	False:     "false",
	And:       "and",
	Or:        "or",
}

var keywords = map[string]Kind{
	"if":    If,
	"while": While,
	"end":   End,
	"true":  True,
	"false": False,
	"and":   And,
	"or":    Or,
}
