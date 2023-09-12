package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"evylang.dev/e/lex"
	"evylang.dev/e/parse"
)

var lexOnly = flag.Bool("l", false, "lexer output")
var help = flag.Bool("help", false, "print help")

func main() {
	flag.Parse()
	if *help {
		printHelp(os.Stdout)
		return
	}
	if len(flag.Args()) > 1 {
		usage()
	}
	input := read()
	tokens := lex.Tokenize(input)
	if *lexOnly {
		for _, tok := range tokens {
			fmt.Println(tok)
		}
		return
	}
	prog, err := parse.ParseProg(tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(prog)
}

func read() string {
	var r io.Reader
	switch len(flag.Args()) {
	case 0:
		r = os.Stdin
	case 1:
		f, err := os.Open(flag.Args()[0])
		handle(err)
		r = f
	default:
		usage()
	}
	b, err := io.ReadAll(r)
	handle(err)
	return string(b)
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintf(w, "usage: [-l] [FILE]\n")
	fmt.Fprintf(w, "  -l: lexer output only \n")
	fmt.Fprintf(w, "  FILE defaults to stdin \n")
}

func usage() {
	printHelp(os.Stderr)
	os.Exit(1)
}
