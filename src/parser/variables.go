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

var Variables = make(map[string]Variable)

func SetVariable(name string, value string, valType string) {
	if valType == "expr" {
		value = strings.TrimPrefix(value, `[`)
		value = strings.TrimSuffix(value, `]`)

		evalParams := make(map[string]any)
		for k, v := range Variables {
			if v.ttype == "bool" {
				evalParams[k] = (v.Value == "true")
			} else if numVal, err := strconv.Atoi(v.Value); err == nil {
				evalParams[k] = numVal
			} else {
				evalParams[k] = v.Value
			}
		}

		expression, err := govaluate.NewEvaluableExpression(value)
		if err != nil {
			fmt.Printf("[Runtime Error]: Invalid expression '%s'\n", value)
			return
		}

		result, err := expression.Evaluate(evalParams)
		if err != nil {
			fmt.Printf("[Runtime Error]: Failed to evaluate '%s': %v\n", value, err)
			return
		}

		value = fmt.Sprint(result)

		switch result.(type) {
		case string:
			valType = "string"
		default:
			valType = "int"
		}
	}

	if valType == "var" {
		existing, ok := Variables[value]
		if !ok {
			fmt.Printf("[Runtime Error]: Variable not found '%s'\n", value)
			return
		}

		value = existing.Value
		valType = existing.ttype
	}

	Variables[name] = Variable{
		Name:  name,
		ttype: valType,
		Value: value,
	}
}
