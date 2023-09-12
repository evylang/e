# e

`e` is a cut down version of [evy](https://github.com/foxygoat/evy).

This repository demonstrates how to build a minimal lexer and parser from
scratch without using any external libraries.

Try the `e` command with:

```sh
go run main.go test.e
```

## Syntax grammar

```
prog = { stmt } .
stmt = decl | assign .

decl  = ident ":" type .
ident = LETTER { LETTER | DIGIT } .
type  = "num" | "string" | "bool" .

assign = ident "=" expr .

expr    = operand | unary_expr |
          binary_expr .
operand = literal | ident | group .
literal = /* e.g. "abc", 1, 2.34, true, false */ .
group   = "(" expr ")" .

unary_expr = UNARY_OP expr .
UNARY_OP   = "-" | "!" .

binary_expr = expr BINARY_OP expr .
BINARY_OP = "*" | "/" | "%" |
            "+" | "-" |
            "<" | "<=" | ">" | ">=" |
            "==" | "!=" |
            "and" |
            "or" .
```

