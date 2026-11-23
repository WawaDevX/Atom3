package lexer

import (
	"fmt"
	"strings"
)

type Token struct {
	Kind  string
	Value string
}

func Tokenize(line string) ([]Token, error) {
	var tokens []Token
	for i := 0; i < len(line); {
		if strings.ContainsRune(" \t\r\n", rune(line[i])) {
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
			tokens = append(tokens, Token{"value", line[start:i]})
		case line[i] == '[':
			i++
			for i < len(line) && line[i] != ']' {
				i++
			}
			if i == len(line) || i == start+1 {
				return nil, fmt.Errorf("invalid expression")
			}
			i++
			tokens = append(tokens, Token{"value", line[start:i]})
		case line[i] == '{' && len(tokens) >= 3 && tokens[0].Value == "var" && tokens[len(tokens)-1].Kind == "=":
			depth := 1
			inString := false
			i++
			for i < len(line) && depth > 0 {
				if line[i] == '"' && (i == 0 || line[i-1] != '\\') {
					inString = !inString
				} else if !inString {
					if line[i] == '{' {
						depth++
					} else if line[i] == '}' {
						depth--
					}
				}
				i++
			}
			if depth != 0 || inString {
				return nil, fmt.Errorf("unterminated array")
			}
			tokens = append(tokens, Token{"value", line[start:i]})
		case IsIdentifierStart(line[i]):
			i++
			for i < len(line) && IsIdentifierPart(line[i]) {
				i++
			}
			if i < len(line) && line[i] == '[' {
				i++
				for i < len(line) && line[i] != ']' {
					i++
				}
				if i == len(line) {
					return nil, fmt.Errorf("unterminated index")
				}
				i++
			}
			tokens = append(tokens, Token{"ident", line[start:i]})
		case line[i] >= '0' && line[i] <= '9' || line[i] == '-' && i+1 < len(line) && line[i+1] >= '0' && line[i+1] <= '9':
			if line[i] == '-' {
				i++
			}
			for i < len(line) && line[i] >= '0' && line[i] <= '9' {
				i++
			}
			tokens = append(tokens, Token{"value", line[start:i]})
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
				tokens = append(tokens, Token{matched, matched})
				continue
			}
			if strings.ContainsRune("{}=(),", rune(line[i])) {
				matched = line[i : i+1]
				i++
				tokens = append(tokens, Token{matched, matched})
				continue
			}
			return nil, fmt.Errorf("unexpected character %q", line[i])
		}
	}
	return tokens, nil
}

func IsIdentifierStart(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func IsIdentifierPart(c byte) bool {
	return IsIdentifierStart(c) || c >= '0' && c <= '9' || c == '.'
}
