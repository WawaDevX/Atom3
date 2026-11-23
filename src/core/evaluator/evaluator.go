package evaluator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

var cache = make(map[string]*vm.Program)

type ArrLenFunc func(string) (int, error)
type Function func(...any) (any, error)

func Compile(source string, environment map[string]any, arrlen ArrLenFunc, functionMaps ...map[string]Function) (*vm.Program, error) {
	var names []string
	for _, functionMap := range functionMaps {
		for name := range functionMap {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	cacheKey := source + "|" + strings.Join(names, ",")
	if cached, exists := cache[cacheKey]; exists {
		return cached, nil
	}

	options := []expr.Option{
		expr.Env(environment),
		expr.Function("arrlen", func(params ...any) (any, error) {
			if len(params) != 1 {
				return nil, fmt.Errorf("arrlen expects 1 argument")
			}
			name, ok := params[0].(string)
			if !ok {
				return nil, fmt.Errorf("arrlen expects an array name")
			}
			return arrlen(name)
		}),
	}
	for _, functionMap := range functionMaps {
		for name, function := range functionMap {
			capturedName := name
			capturedFunction := function
			options = append(options, expr.Function(capturedName, func(params ...any) (any, error) {
				return capturedFunction(params...)
			}))
		}
	}

	compiled, err := expr.Compile(source, options...)
	if err == nil {
		cache[cacheKey] = compiled
	}
	return compiled, err
}

func Evaluate(source string, environment map[string]any, arrlen ArrLenFunc, functionMaps ...map[string]Function) (any, error) {
	program, err := Compile(source, environment, arrlen, functionMaps...)
	if err != nil {
		return nil, err
	}
	return expr.Run(program, environment)
}
