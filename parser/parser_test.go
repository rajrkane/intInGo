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

  // invalid test case
  // input := `
  //   let x 5;
  //   let = 10;
  //   let 8383882;
  // `

  // TODO: look at more parser tests

  // initialize new lexer, parser
  l := lexer.New(input)
  p := New(l)

  program := p.ParseProgram()
  checkParserErrors(t, p)
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

func TestReturnStatements(t *testing.T) {
  input := `
    return 5;
    return 10;
    return 999293;
  `

  l := lexer.New(input)
  p := New(l)

  program := p.ParseProgram()
  checkParserErrors(t, p)
  if len(program.Statements) != 3 {
    t.Fatalf("program expects 3 statements, got=%d", len(program.Statements))
  }

  for _, stmt := range program.Statements {
    returnStmt, ok := stmt.(*ast.ReturnStatement)
    if !ok {
      t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
      continue
    }
    if returnStmt.TokenLiteral() != "return" {
      t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
    }
  }
}

func TestIdentifierExpression(t *testing.T) {
  input := "foobar;"

  l := lexer.New(input)
  p := New(l)

  program := p.ParseProgram()
  checkParserErrors(t, p)
  if len(program.Statements) != 1 {
    t.Fatalf("program expects 1 statement, got=%d", len(program.Statements))
  }

  // only 1 statement, so no need to loop
  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Errorf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
  }

  ident, ok := stmt.Expression.(*ast.Identifier)
  if !ok {
    t.Fatalf("expression not identifier, got=%d", stmt.Expression)
  }
  if ident.Value != "foobar" {
    t.Errorf("identifier value not %s, got %s", "foobar", ident.Value)
  }
  if ident.TokenLiteral() != "foobar" {
    t.Errorf("identifier token literal not %s, got %s", "foobar", ident.TokenLiteral())
  }
}

// print any parser errors
func checkParserErrors(t *testing.T, p *Parser) {
  errors := p.Errors()
  if len(errors) == 0 {
    return
  }

  t.Errorf("parser has %d errors", len(errors))
  for _, msg := range errors {
    t.Errorf("parser error: %q", msg)
  }
  t.FailNow()
}
