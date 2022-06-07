// parser/parser_test.go

package parser

import (
  "testing"
  "intInGo/ast"
  "intInGo/lexer"
)

func TestLetStatements(t *testing.T) {

  // test case
  input := `
    let x = 5;
    let y = 10;
    let foobar = 8822813;
  `

  // initialize new lexer, parser
  l := lexer.New(input)
  p := New(l)

  program := p.ParseProgram()
  if program == nil {
    t.Fatalf("ParseProgram() returned nil")
  }
  if len(program.Statements) != 3 { // hard-coded to 3 statements for now
    t.Fatalf("program expects 3 statements, got=%d", len(program.Statements))
  }

  tests := []struct {
    expectedIdentifier string
  }{
    {"x"},
    {"y"},
    {"foobar"},
  }

  for i, tt := range tests {
    stmt := program.Statements[i]
    // test each let statement
    if !testLetStatement(t, stmt, tt.expectedIdentifier) {
      return
    }
  }
}

// check as many fields of an AST node as possible
func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
  // test if let statement
  if stmt.TokenLiteral() != "let" {
    t.Errorf("stmt.TokenLiteral not 'let', got=%q", stmt.TokenLiteral())
    return false
  }
  letStmt, ok := stmt.(*ast.LetStatement)
  if !ok {
    t.Errorf("stmt not *ast.LetStatement, got=%T", stmt)
    return false
  }

  // test identifier of the binding
  if letStmt.Name.Value != name {
    t.Errorf("letStmt.Name.Value not '%s', got=%s", name, letStmt.Name.Value)
    return false
  }
  if letStmt.Name.TokenLiteral() != name {
    t.Errorf("letStmt.Name.TokenLiteral() not %s, got=%s", name, letStmt.Name.TokenLiteral())
    return false
  }

  return true
}
