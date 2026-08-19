package parser

import (
	"os"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var atomLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Expr", Pattern: `\[[^\]]+\]`},
	{Name: "String", Pattern: `"[^"]*"`},
	{Name: "Int", Pattern: `\d+`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Punct", Pattern: `[=(),]`},
	{Name: "whitespace", Pattern: `\s+`},
})

type VarDecl struct {
	Name  string `@Ident "="`
	Value string `@Expr | @String | @Int | @Ident`
}

type FunctionCall struct {
	Name string   `@Ident`
	Args []string `"(" [ (@String | @Ident) ( "," (@String | @Ident) )* ] ")"`
}

type Statement struct {
	VarCall  *VarDecl      `  "var" @@`
	FuncCall *FunctionCall `| @@`
}

var Parser = participle.MustBuild[Statement](participle.Lexer(atomLexer))

func isNumber(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func RunLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	stmt, err := Parser.ParseString("", line)
	if err != nil || stmt == nil {
		return
	}

	if stmt.VarCall != nil {
		val := stmt.VarCall.Value
		varType := "var"

		if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			varType = "expr"
		} else if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
			varType = "string"
		} else if isNumber(val) {
			varType = "int"
		}

		SetVariable(stmt.VarCall.Name, val, varType)
		return
	}

	if stmt.FuncCall != nil {
		call := stmt.FuncCall

		if call.Name == "print" {
			for _, arg := range call.Args {
				if strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`) {
					Print(arg, "string")
				} else if isNumber(arg) {
					Print(arg, "int")
				} else {
					Print(arg, "var")
				}
			}
		} else if call.Name == "exit" {
			os.Exit(0)
		}
	}
}
