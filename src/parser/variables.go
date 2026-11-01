package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

type Variable struct {
	Name  string
	ttype string
	Value string
}

type arrayElement struct {
	ttype string
	Value string
}

type array struct {
	Name     string
	Elements []arrayElement
}

var Variables = make(map[string]Variable)
var Arrays = make(map[string]array)
var exprCache = make(map[string]*govaluate.EvaluableExpression)
var evalParams = make(map[string]any)

func splitArray(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	var elements []string
	start := 0
	braceDepth := 0
	bracketDepth := 0
	inString := false
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '"':
			if i == 0 || value[i-1] != '\\' {
				inString = !inString
			}
		case '{':
			if !inString {
				braceDepth++
			}
		case '}':
			if !inString {
				braceDepth--
			}
		case '[':
			if !inString {
				bracketDepth++
			}
		case ']':
			if !inString {
				bracketDepth--
			}
		case ',':
			if !inString && braceDepth == 0 && bracketDepth == 0 {
				elements = append(elements, strings.TrimSpace(value[start:i]))
				start = i + 1
			}
		}
	}
	elements = append(elements, strings.TrimSpace(value[start:]))
	return elements
}

func SetVariable(name string, value string, valType string) {
	if valType == "expr" {
		value = strings.TrimPrefix(value, `[`)
		value = strings.TrimSuffix(value, `]`)

		expression, exists := exprCache[value]
		if !exists {
			var err error
			expression, err = govaluate.NewEvaluableExpression(value)
			if err != nil {
				fmt.Printf("[Runtime Error]: Invalid expression '%s'\n", value)
				return
			}
			exprCache[value] = expression
		}

		result, err := expression.Evaluate(evalParams)
		if err != nil {
			fmt.Printf("[Runtime Error]: Failed to evaluate '%s': %v\n", value, err)
			return
		}

		if f, ok := result.(float64); ok {
			value = strconv.FormatFloat(f, 'f', -1, 64)
		} else {
			value = fmt.Sprint(result)
		}

		switch result.(type) {
		case bool:
			valType = "bool"
		case int, int64, float64:
			valType = "int"
		default:
			valType = "string"
		}
	}

	if valType == "var" {
		existing, ok := Variables[value]
		if !ok {
			fmt.Printf("[Runtime Error]: Variable not found or has not been defined yet!  '%s'\n", value)
			return
		}

		value = existing.Value
		valType = existing.ttype
	}

	if valType == "int" {
		if numVal, err := strconv.Atoi(value); err == nil {
			evalParams[name] = numVal
		} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			evalParams[name] = floatVal
		} else {
			evalParams[name] = value
		}
	} else if valType == "bool" {
		evalParams[name] = (value == "true")
	} else if valType == "array" {
		arrayContents := strings.TrimSuffix(strings.TrimPrefix(value, `{`), `}`)
		elements := splitArray(arrayContents)
		var arrayElements []arrayElement
		for _, elem := range elements {
			arrayElements = append(arrayElements, arrayElement{
				ttype: DetermineType(elem),
				Value: elem,
			})
		}
		Arrays[name] = array{
			Name:     name,
			Elements: arrayElements,
		}
		delete(evalParams, name)
	} else {
		delete(Arrays, name)
		value = strings.Trim(value, `" `)
		evalParams[name] = value
	}

	if existing, ok := Variables[name]; ok {
		existing.ttype = valType
		existing.Value = value
		Variables[name] = existing
	} else {
		Variables[name] = Variable{
			Name:  name,
			ttype: valType,
			Value: value,
		}
	}

}
