package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var skipBlock = false

var atomLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Expr", Pattern: `\[[^\]]+\]`},
	{Name: "String", Pattern: `"[^"]*"`},
	{Name: "Int", Pattern: `\d+`},
	{Name: "LBrace", Pattern: `\{`},
	{Name: "RBrace", Pattern: `\}`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Punct", Pattern: `==|\?=|\!=|<=|>=|[=(),]`}, // Added ?= here
	{Name: "whitespace", Pattern: `\s+`},
})

type IfStmt struct {
	Condition string `@Expr "{"`
}

type VarDecl struct {
	Name  string `@Ident`
	Value string `"=" ( @Expr | @String | @Int | @Ident )`
}

type FunctionCall struct {
	Name string   `@Ident`
	Args []string `"(" [ (@String | @Ident) ( "," (@String | @Ident) )* ] ")"`
}

type Statement struct {
	IfCall   *IfStmt       `  "if" @@`
	VarCall  *VarDecl      `| "var" @@`
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

	if line == "}" {
		skipBlock = false
		return
	}

	if skipBlock {
		return
	}

	stmt, err := Parser.ParseString("", line)
	if err != nil {
		fmt.Printf("[Parse Error] on line '%s': %v\n", line, err)
		return
	}
	if stmt == nil {
		return
	}

	// if
	if stmt.IfCall != nil {
		cond := strings.TrimSuffix(strings.TrimPrefix(stmt.IfCall.Condition, "["), "]")

		// Convert Atom 3.0's custom ?= operator to != for govaluate
		cond = strings.ReplaceAll(cond, "?=", "!=")

		evalParams := make(map[string]any)
		for k, v := range Variables {
			cleanVal := strings.Trim(v.Value, `" `)

			if numVal, err := strconv.Atoi(cleanVal); err == nil {
				evalParams[k] = numVal
			} else if floatVal, err := strconv.ParseFloat(cleanVal, 64); err == nil {
				evalParams[k] = floatVal
			} else {
				evalParams[k] = cleanVal
			}
		}

		expr, err := govaluate.NewEvaluableExpression(cond)
		if err != nil {
			skipBlock = true
			return
		}

		result, err := expr.Evaluate(evalParams)
		if err != nil {
			skipBlock = true
			return
		}

		boolResult, isBool := result.(bool)
		if !isBool || !boolResult {
			skipBlock = true
		}
		return
	}

	// variables
	if stmt.VarCall != nil {
		val := stmt.VarCall.Value
		varName := strings.TrimSpace(stmt.VarCall.Name)
		varType := "ident"

		if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			varType = "expr"
		} else if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
			varType = "string"
		} else if isNumber(val) {
			varType = "int"
		} else {
			varType = "var"
		}

		SetVariable(varName, val, varType)
		return
	}

	// Functions
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
