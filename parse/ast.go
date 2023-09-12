package parse

import (
	"fmt"
	"strings"

	"evylang.dev/e/lex"
)

type Node interface {
	String() string
	Type() Type
}

type Type string

const (
	NumType    Type = "num"
	StringType Type = "string"
	BoolType   Type = "bool"
	NoneType   Type = "none"
)

func (t Type) String() string {
	return string(t)
}

type Prog struct {
	Stmts []Node
}

func (p *Prog) String() string {
	sb := &strings.Builder{}
	sb.WriteString("PROG {\n")
	for _, s := range p.Stmts {
		sb.WriteString("\t" + s.String() + "\n")
	}
	sb.WriteString("}")
	return sb.String()
}

func (p *Prog) Type() Type { return NoneType }

type Assign struct {
	Target *Var
	Val    Node
}

func (a *Assign) Type() Type     { return a.Target.Type() }
func (a *Assign) String() string { return "ASSIGN " + a.Target.String() + " = " + a.Val.String() }

type Var struct {
	Name string
	T    Type
}

func (v *Var) Type() Type     { return v.T }
func (v *Var) String() string { return v.Name }

type Decl struct {
	Var *Var
}

func (d *Decl) Type() Type     { return d.Var.T }
func (d *Decl) String() string { return "DECL   " + d.Var.Name + ":" + string(d.Var.T) }

type NumLit struct {
	Val float64
}

func (n *NumLit) Type() Type     { return NumType }
func (n *NumLit) String() string { return fmt.Sprintf("%v", n.Val) }

type BoolLit struct {
	Val bool
}

func (n *BoolLit) Type() Type     { return BoolType }
func (n *BoolLit) String() string { return fmt.Sprintf("%v", n.Val) }

type StringLit struct {
	Val string
}

func (s *StringLit) Type() Type     { return StringType }
func (s *StringLit) String() string { return fmt.Sprintf("%q", s.Val) }

type BinaryExpr struct {
	Left  Node
	Op    lex.Kind
	Right Node
}

func (b *BinaryExpr) Type() Type {
	if b.Op.IsComparisonOp() {
		return BoolType
	}
	return b.Left.Type()
}
func (b *BinaryExpr) String() string { return fmt.Sprintf("(%v %s %v)", b.Left, b.Op.Repr(), b.Right) }

type UnaryExpr struct {
	Op    lex.Kind
	Right Node
}

func (u *UnaryExpr) Type() Type     { return u.Right.Type() }
func (u *UnaryExpr) String() string { return fmt.Sprintf("(%s %v)", u.Op.Repr(), u.Right) }

type GroupExpr struct {
	Expr Node
}

func (g *GroupExpr) Type() Type     { return g.Expr.Type() }
func (g *GroupExpr) String() string { return fmt.Sprintf("%v", g.Expr) }
