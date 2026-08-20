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
var exprCache = make(map[string]*govaluate.EvaluableExpression)
var evalParams = make(map[string]any)

func updateEvalParams() {
	for k, v := range Variables {
		if v.ttype == "int" {
			if numVal, err := strconv.Atoi(v.Value); err == nil {
				evalParams[k] = numVal
				continue
			}
		} else if v.ttype == "bool" {
			evalParams[k] = (v.Value == "true")
			continue
		}
		evalParams[k] = strings.Trim(v.Value, `" `)
	}
}

func SetVariable(name string, value string, valType string) {
	if valType == "expr" {
		value = strings.TrimPrefix(value, `[`)
		value = strings.TrimSuffix(value, `]`)

		updateEvalParams()

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
	} else {
		evalParams[name] = strings.Trim(value, `" `)
	}
}
