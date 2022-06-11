// parser/parser_test.go

package parser

import (
  "testing"
  "intInGo/ast"
  "intInGo/lexer"
	"fmt"
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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program expects 1 statement, got=%d", len(program.Statements))
	}

	// only 1 statement
	// assert that the first statement is an expression statement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	// expect an integer literal
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not integer literal, got=%d", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal value not %d, got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal token literal not %s, got %s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input					string
		operator			string
		integerValue	int64
	}{
		{"!15;", "!", 15},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program expects 1 statement, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement not prefix expression, got=%d", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input					string
		leftValue			int64
		operator			string
		rightValue	int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program expects 1 statement, got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ExpressionStatement, got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement not infix expression, got=%d", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
		integ, ok := il.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
			return false
		}
		if integ.Value != value {
			t.Errorf("integ.Value not %d, got=%d", value, integ.Value)
		}
		if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
			t.Errorf("integ.TokenLiteral not %d, got=%s", value, integ.TokenLiteral())
			return false
		}
		return true
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
