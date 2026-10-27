package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
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

type IfStmt struct {
	Condition string
}

type VarDecl struct {
	Name    string
	FuncVal *FunctionCall
	Value   string
}

type FunctionCall struct {
	Name string
	Args []string
}

type ReturnStmt struct {
	Value string
}

type Statement struct {
	IfCall     *IfStmt
	VarCall    *VarDecl
	ReturnCall *ReturnStmt
	FuncCall   *FunctionCall
}

type statementToken struct {
	kind  string
	value string
}

func tokenizeStatement(line string) ([]statementToken, error) {
	var tokens []statementToken
	for i := 0; i < len(line); {
		if line[i] == ' ' || line[i] == '\t' || line[i] == '\r' || line[i] == '\n' {
			i++
			continue
		}
		if i+1 < len(line) && line[i:i+2] == "//" {
			break
		}

		start := i
		switch {
		case line[i] == '"':
			i++
			for i < len(line) && line[i] != '"' {
				i++
			}
			if i == len(line) {
				return nil, fmt.Errorf("unterminated string")
			}
			i++
			tokens = append(tokens, statementToken{kind: "value", value: line[start:i]})
		case line[i] == '[':
			i++
			for i < len(line) && line[i] != ']' {
				i++
			}
			if i == len(line) || i == start+1 {
				return nil, fmt.Errorf("invalid expression")
			}
			i++
			tokens = append(tokens, statementToken{kind: "value", value: line[start:i]})
		case isIdentifierStart(line[i]):
			i++
			for i < len(line) && isIdentifierPart(line[i]) {
				i++
			}
			tokens = append(tokens, statementToken{kind: "ident", value: line[start:i]})
		case line[i] >= '0' && line[i] <= '9' || line[i] == '-' && i+1 < len(line) && line[i+1] >= '0' && line[i+1] <= '9':
			if line[i] == '-' {
				i++
			}
			for i < len(line) && line[i] >= '0' && line[i] <= '9' {
				i++
			}
			tokens = append(tokens, statementToken{kind: "value", value: line[start:i]})
		default:
			matched := ""
			for _, op := range []string{"==", "?=", "!=", "<=", ">="} {
				if strings.HasPrefix(line[i:], op) {
					matched = op
					break
				}
			}
			if matched != "" {
				i += len(matched)
				tokens = append(tokens, statementToken{kind: matched, value: matched})
				continue
			}
			if strings.ContainsRune("{}=(),", rune(line[i])) {
				matched = line[i : i+1]
				i++
				tokens = append(tokens, statementToken{kind: matched, value: matched})
				continue
			}
			return nil, fmt.Errorf("unexpected character %q", line[i])
		}
	}
	return tokens, nil
}

func isIdentifierStart(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isIdentifierPart(c byte) bool {
	return isIdentifierStart(c) || c >= '0' && c <= '9'
}

type statementParser struct {
	tokens []statementToken
	pos    int
}

func (p *statementParser) accept(kind string) (statementToken, bool) {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].kind != kind {
		return statementToken{}, false
	}
	token := p.tokens[p.pos]
	p.pos++
	return token, true
}

func (p *statementParser) expect(kind string) (statementToken, error) {
	if token, ok := p.accept(kind); ok {
		return token, nil
	}
	if p.pos >= len(p.tokens) {
		return statementToken{}, fmt.Errorf("expected %q at end of statement", kind)
	}
	return statementToken{}, fmt.Errorf("expected %q, got %q", kind, p.tokens[p.pos].value)
}

func (p *statementParser) parseValue() (string, error) {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].kind != "value" && p.tokens[p.pos].kind != "ident" {
		return "", fmt.Errorf("expected a literal, expression, or identifier")
	}
	value := p.tokens[p.pos].value
	p.pos++
	return value, nil
}

func (p *statementParser) parseFunctionCall() (*FunctionCall, error) {
	name, err := p.expect("ident")
	if err != nil {
		return nil, err
	}
	if _, err = p.expect("("); err != nil {
		return nil, err
	}
	call := &FunctionCall{Name: name.value}
	if _, ok := p.accept(")"); ok {
		return call, nil
	}
	for {
		arg, parseErr := p.parseValue()
		if parseErr != nil {
			return nil, parseErr
		}
		call.Args = append(call.Args, arg)
		if _, ok := p.accept(")"); ok {
			return call, nil
		}
		if _, err = p.expect(","); err != nil {
			return nil, err
		}
	}
}

func parseStatement(line string) (*Statement, error) {
	tokens, err := tokenizeStatement(line)
	if err != nil {
		return nil, err
	}
	p := &statementParser{tokens: tokens}
	stmt := &Statement{}

	if keyword, ok := p.accept("ident"); ok {
		switch keyword.value {
		case "if":
			condition, parseErr := p.parseValue()
			if parseErr != nil || !strings.HasPrefix(condition, "[") {
				if parseErr == nil {
					parseErr = fmt.Errorf("if condition must be an expression")
				}
				return nil, parseErr
			}
			if _, err = p.expect("{"); err != nil {
				return nil, err
			}
			stmt.IfCall = &IfStmt{Condition: condition}
		case "var":
			name, parseErr := p.expect("ident")
			if parseErr != nil {
				return nil, parseErr
			}
			if _, err = p.expect("="); err != nil {
				return nil, err
			}
			decl := &VarDecl{Name: name.value}
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos].kind == "ident" && p.tokens[p.pos+1].kind == "(" {
				decl.FuncVal, err = p.parseFunctionCall()
			} else {
				decl.Value, err = p.parseValue()
			}
			if err != nil {
				return nil, err
			}
			stmt.VarCall = decl
		case "return":
			ret := &ReturnStmt{}
			if p.pos < len(p.tokens) {
				ret.Value, err = p.parseValue()
				if err != nil {
					return nil, err
				}
			}
			stmt.ReturnCall = ret
		default:
			p.pos = 0
			stmt.FuncCall, err = p.parseFunctionCall()
		}
	} else {
		return nil, fmt.Errorf("expected statement keyword or function call")
	}
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.tokens) {
		return nil, fmt.Errorf("unexpected token %q", p.tokens[p.pos].value)
	}
	return stmt, nil
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func evaluateCondition(condStr string) bool {
	cond := strings.TrimSuffix(strings.TrimPrefix(condStr, "["), "]")
	cond = strings.ReplaceAll(cond, "?=", "!=")

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

	stmt, err := parseStatement(line)
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

		if stmt.VarCall.FuncVal != nil {
			result := ExecuteFunction(stmt.VarCall.FuncVal.Name, stmt.VarCall.FuncVal.Args)
			if !result.returned {
				return
			}
			val = result.value
			varType = result.ttype
		}

		if varType == "ident" {
			if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
				varType = "expr"
			} else if val == "true" || val == "false" {
				varType = "bool"
			} else if strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`) {
				varType = "string"
			} else if isNumeric(val) {
				varType = "int"
			} else {
				varType = "var"
			}
		}

		SetVariable(varName, val, varType)
		return
	}

	if stmt.ReturnCall != nil {
		if len(functionReturnStack) > 0 {
			frame := functionReturnStack[len(functionReturnStack)-1]
			frame.value = stmt.ReturnCall.Value
			frame.returned = true
		}
		return
	}

	// Functions
	if stmt.FuncCall != nil {
		ExecuteFunction(stmt.FuncCall.Name, stmt.FuncCall.Args)
	}
}
