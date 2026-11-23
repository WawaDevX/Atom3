package parser

import (
	"fmt"
	"strconv"
	"strings"

	"Atom3/src/core/evaluator"
	atomlexer "Atom3/src/core/lexer"
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
var currentFuncDefaults map[string]string
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

type statementParser struct {
	tokens []atomlexer.Token
	pos    int
}

func (p *statementParser) accept(kind string) (atomlexer.Token, bool) {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Kind != kind {
		return atomlexer.Token{}, false
	}
	token := p.tokens[p.pos]
	p.pos++
	return token, true
}

func (p *statementParser) expect(kind string) (atomlexer.Token, error) {
	if token, ok := p.accept(kind); ok {
		return token, nil
	}
	if p.pos >= len(p.tokens) {
		return atomlexer.Token{}, fmt.Errorf("expected %q at end of statement", kind)
	}
	return atomlexer.Token{}, fmt.Errorf("expected %q, got %q", kind, p.tokens[p.pos].Value)
}

func (p *statementParser) parseValue() (string, error) {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Kind != "value" && p.tokens[p.pos].Kind != "ident" {
		RaiseException("expected a literal, expression, or identifier", ErrSyntax)
	}
	value := p.tokens[p.pos].Value
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
	call := &FunctionCall{Name: name.Value}
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
	tokens, err := atomlexer.Tokenize(line)
	if err != nil {
		return nil, err
	}
	p := &statementParser{tokens: tokens}
	stmt := &Statement{}

	if keyword, ok := p.accept("ident"); ok {
		switch keyword.Value {
		case "if":
			condition, parseErr := p.parseValue()
			if parseErr != nil || !strings.HasPrefix(condition, "[") {
				if parseErr == nil {
					parseErr = RaiseException("if condition must be an expression", ErrSyntax)
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
			decl := &VarDecl{Name: name.Value}
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos].Kind == "ident" && p.tokens[p.pos+1].Kind == "(" {
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
		return nil, RaiseException("expected statement keyword or function call", ErrSyntax)
	}
	if err != nil {
		return nil, err
	}
	if p.pos != len(p.tokens) {
		return nil, fmt.Errorf("unexpected token %q", p.tokens[p.pos].Value)
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

	result, err := evaluator.Evaluate(cond, evalParams, arrLen, expressionFunctions())
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
			name := currentFuncName
			if currentModule != "" {
				name = currentModule + "." + name
			}
			if _, exists := UserFunctions[name]; exists {
				RaiseException(
					fmt.Sprintf("Function '%s' is already defined", name),
					ErrRuntime,
				)
			} else {
				UserFunctions[name] = UserFunc{
					Params:   currentFuncParams,
					Defaults: currentFuncDefaults,
					Body:     currentFuncLines,
					Module:   currentModule,
				}
			}
			currentFuncLines = nil
			currentFuncParams = nil
			currentFuncDefaults = nil
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
			currentFuncDefaults = make(map[string]string)
			if paramStr != "" {
				rawParams := strings.Split(paramStr, ",")
				for _, raw := range rawParams {
					param := strings.TrimSpace(raw)
					if name, defaultValue, found := strings.Cut(param, "="); found {
						name = strings.TrimSpace(name)
						defaultValue = strings.TrimSpace(defaultValue)
						if name == "" || defaultValue == "" {
							RaiseException(fmt.Sprintf("Invalid optional parameter '%s'", param), ErrSyntax)
							return
						}
						currentFuncParams = append(currentFuncParams, name)
						currentFuncDefaults[name] = defaultValue
					} else if param != "" {
						currentFuncParams = append(currentFuncParams, param)
					}
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
			RaiseException(
				fmt.Sprintf("'else' without matching 'if' on line '%s'", line),
				ErrParse,
			)
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

	if strings.HasPrefix(line, "import ") {
		if currentlyActive() {
			ImportFile(strings.TrimSpace(strings.TrimPrefix(line, "import ")))
		}
		return
	}

	stmt, err := parseStatement(line)
	if err != nil {
		RaiseException(
			fmt.Sprintf("on line '%s': %v", line, err),
			ErrParse,
		)
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
			} else if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
				varType = "array"
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
		if len(functionReturnStack) == 0 {
			RaiseException("'return' cannot be used outside a function", ErrInvalidOperation)
			return
		}
		frame := functionReturnStack[len(functionReturnStack)-1]
		frame.value = ResolveValue(stmt.ReturnCall.Value, true)
		if variable, exists := Variables[stmt.ReturnCall.Value]; exists {
			frame.ttype = variable.ttype
		} else {
			frame.ttype = DetermineType(frame.value)
		}
		frame.returned = true
		return
	}

	// Functions
	if stmt.FuncCall != nil {
		ExecuteFunction(stmt.FuncCall.Name, stmt.FuncCall.Args)
	}
}
