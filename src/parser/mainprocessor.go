package parser

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type blockFrame struct {
	active             bool
	taken              bool
	parentActiveAtOpen bool
}

var blockStack []blockFrame

type elseInfo struct {
	parentActive bool
	ifTaken      bool
}

var lastIfClosed *elseInfo

func currentlyActive() bool {
	if len(blockStack) == 0 {
		return true
	}
	return blockStack[len(blockStack)-1].active
}

var whileDepth = 0

// Loop Buffering
var isBufferingWhile = false
var whileCond = ""
var whileLines []string

var isBufferingFunc = false
var currentFuncName = ""
var currentFuncParams []string
var currentFuncLines []string

var atomLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comment", Pattern: `//.*`},
	{Name: "Expr", Pattern: `\[[^\]]+\]`},
	{Name: "String", Pattern: `"[^"]*"`},
	{Name: "Bool", Pattern: `\b(true|false)\b`},
	{Name: "Int", Pattern: `\d+`},
	{Name: "LBrace", Pattern: `\{`},
	{Name: "RBrace", Pattern: `\}`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Punct", Pattern: `==|\?=|\!=|<=|>=|[=(),]`},
	{Name: "whitespace", Pattern: `\s+`},
})

type IfStmt struct {
	Condition string `@Expr "{"`
}

type ElseStmt struct {
	Body string `"else" "{"`
}

type VarDecl struct {
	Name    string        `@Ident "="`
	FuncVal *FunctionCall `( @@`
	Value   string        `| @Expr | @String | @Bool | @Int | @Ident )`
}

type FunctionCall struct {
	Name string   `@Ident`
	Args []string `"(" [ (@String | @Bool | @Int | @Expr | @Ident) ( "," (@String | @Bool | @Int | @Expr | @Ident) )* ] ")"`
}

type ReturnStmt struct {
	Value string `"return" [ @Expr | @String | @Bool | @Int | @Ident ]`
}

type Statement struct {
	IfCall     *IfStmt       `  "if" @@`
	ElseCall   *ElseStmt     `| @@`
	VarCall    *VarDecl      `| "var" @@`
	ReturnCall *ReturnStmt   `| @@`
	FuncCall   *FunctionCall `| @@`
}

var Parser = participle.MustBuild[Statement](participle.Lexer(atomLexer), participle.Elide("Comment"))

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

func evaluateCondition(condStr string) bool {
	cond := strings.TrimSuffix(strings.TrimPrefix(condStr, "["), "]")
	cond = strings.ReplaceAll(cond, "?=", "!=")

	updateEvalParams()

	expr, exists := exprCache[cond]
	if !exists {
		var err error
		expr, err = govaluate.NewEvaluableExpression(cond)
		if err != nil {
			return false
		}
		exprCache[cond] = expr
	}

	result, err := expr.Evaluate(evalParams)
	if err != nil {
		return false
	}

	boolResult, isBool := result.(bool)
	return isBool && boolResult
}

func RunLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "//") {
		return
	}

	if isBufferingFunc {
		if line == "}" {
			isBufferingFunc = false
			UserFunctions[currentFuncName] = UserFunc{
				Params: currentFuncParams,
				Body:   currentFuncLines,
			}
			currentFuncLines = nil
			currentFuncParams = nil
			currentFuncName = ""
			return
		}
		currentFuncLines = append(currentFuncLines, line)
		return
	}

	if strings.HasPrefix(line, "func ") {
		header := strings.TrimPrefix(line, "func ")
		header = strings.TrimSuffix(header, "{")
		header = strings.TrimSpace(header)

		parts := strings.SplitN(header, "(", 2)
		if len(parts) == 2 {
			currentFuncName = strings.TrimSpace(parts[0])
			paramStr := strings.TrimSuffix(parts[1], ")")
			paramStr = strings.TrimSpace(paramStr)

			currentFuncParams = nil
			if paramStr != "" {
				rawParams := strings.Split(paramStr, ",")
				for _, p := range rawParams {
					currentFuncParams = append(currentFuncParams, strings.TrimSpace(p))
				}
			}

			isBufferingFunc = true
			currentFuncLines = []string{}
			return
		}
	}

	// while buffer check
	if isBufferingWhile {
		if line == "}" && whileDepth == 0 {
			isBufferingWhile = false

			// Execute the loop!
			cond := whileCond
			body := whileLines
			whileLines = nil
			whileCond = ""

			for evaluateCondition(cond) {
				for _, loopLine := range body {
					RunLine(loopLine)
				}
			}
			return
		}
		whileDepth += strings.Count(line, "{") - strings.Count(line, "}")
		if whileDepth < 0 {
			whileDepth = 0
		}
		whileLines = append(whileLines, line)
		return
	}

	// while
	if strings.HasPrefix(line, "while") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) > 1 {
			condPart := strings.TrimSpace(parts[1])
			condPart = strings.TrimSuffix(condPart, "{")
			whileCond = strings.TrimSpace(condPart)
			isBufferingWhile = true
			whileDepth = 0
			whileLines = []string{}
			return
		}
	}

	if strings.HasPrefix(line, "}") && line != "}" {
		rest := strings.TrimSpace(line[1:])
		RunLine("}")
		if rest != "" {
			RunLine(rest)
		}
		return
	}

	if !strings.HasPrefix(line, "else") {
		lastIfClosed = nil
	}

	if line == "}" {
		if len(blockStack) > 0 {
			closed := blockStack[len(blockStack)-1]
			blockStack = blockStack[:len(blockStack)-1]
			lastIfClosed = &elseInfo{
				parentActive: closed.parentActiveAtOpen,
				ifTaken:      closed.taken,
			}
		}
		return
	}

	if strings.HasPrefix(line, "else") {
		if lastIfClosed == nil {
			fmt.Printf("[Parse Error] 'else' without matching 'if' on line '%s'\n", line)
			return
		}
		parentActive := lastIfClosed.parentActive
		elseActive := parentActive && !lastIfClosed.ifTaken
		blockStack = append(blockStack, blockFrame{
			active:             elseActive,
			taken:              elseActive,
			parentActiveAtOpen: parentActive,
		})
		lastIfClosed = nil
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
		parentActive := currentlyActive()
		frame := blockFrame{parentActiveAtOpen: parentActive}
		if parentActive {
			cond := evaluateCondition(stmt.IfCall.Condition)
			frame.active = cond
			frame.taken = cond
		}
		blockStack = append(blockStack, frame)
		return
	}

	if !currentlyActive() {
		return
	}

	// variables
	if stmt.VarCall != nil {
		val := stmt.VarCall.Value
		varName := strings.TrimSpace(stmt.VarCall.Name)
		varType := "ident"

		if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			varType = "expr"
		} else if val == "true" || val == "false" {
			varType = "bool"
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
		ExecuteFunction(stmt.FuncCall.Name, stmt.FuncCall.Args)
	}
}
