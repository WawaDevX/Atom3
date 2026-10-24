package parser

import "testing"

func TestParseStatement(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		check func(*testing.T, *Statement)
	}{
		{
			name: "variable declaration",
			line: `var title="hello world"`,
			check: func(t *testing.T, stmt *Statement) {
				if stmt.VarCall == nil || stmt.VarCall.Name != "title" || stmt.VarCall.Value != `"hello world"` {
					t.Fatalf("unexpected declaration: %#v", stmt.VarCall)
				}
			},
		},
		{
			name: "function call with arguments",
			line: `print("hello", true, 42, [score]) // comment`,
			check: func(t *testing.T, stmt *Statement) {
				want := []string{`"hello"`, "true", "42", "[score]"}
				if stmt.FuncCall == nil || stmt.FuncCall.Name != "print" || len(stmt.FuncCall.Args) != len(want) {
					t.Fatalf("unexpected call: %#v", stmt.FuncCall)
				}
				for i := range want {
					if stmt.FuncCall.Args[i] != want[i] {
						t.Errorf("argument %d = %q, want %q", i, stmt.FuncCall.Args[i], want[i])
					}
				}
			},
		},
		{
			name: "if condition",
			line: `if [score >= 10] {`,
			check: func(t *testing.T, stmt *Statement) {
				if stmt.IfCall == nil || stmt.IfCall.Condition != "[score >= 10]" {
					t.Fatalf("unexpected if statement: %#v", stmt.IfCall)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parseStatement(tt.line)
			if err != nil {
				t.Fatalf("parseStatement() error = %v", err)
			}
			tt.check(t, stmt)
		})
	}
}

func TestParseStatementRejectsMalformedInput(t *testing.T) {
	for _, line := range []string{`print("unterminated)`, `var value=`, `print(1,)`, `if score {`} {
		if _, err := parseStatement(line); err == nil {
			t.Errorf("parseStatement(%q) expected error", line)
		}
	}
}
