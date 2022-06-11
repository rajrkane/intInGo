// parser/parser.go

package parser

import (
  "intInGo/ast"
  "intInGo/lexer"
  "intInGo/token"
  "fmt"
	"strconv"
)

type Parser struct {
  l *lexer.Lexer        // ptr to instance of lexer
  errors []string

  curToken  token.Token // point to current token
  peekToken token.Token // point to next token

  // used to check if a prefix or infix fn is associated with curToken.Type
  prefixParseFns  map[token.TokenType]prefixParseFn
  infixParseFns   map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
  p := &Parser{
    l: l,
    errors: []string{},
  }

	// initialize prefix parsing functions map
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

  // read two tokens
  p.nextToken()
  p.nextToken()

  return p
}

// advance both current and next tokens
func (p *Parser) nextToken() {
  p.curToken = p.peekToken
  p.peekToken = p.l.NextToken() // request new token from lexer
}

// entry point to recursive descent parser
func (p *Parser) ParseProgram() *ast.Program {
  // construct root of AST
  program := &ast.Program{}

  program.Statements = []ast.Statement{}

  // build child nodes by iterating over every input token
  for !p.curTokenIs(token.EOF) {
    stmt := p.parseStatement()
    if stmt != nil {
      // add parsed statement to Statements slice of root
      program.Statements = append(program.Statements, stmt)
    }
    // advance current and next tokens
    p.nextToken()
  }

  // return once nothing left to parse
  return program
}

func (p *Parser) parseStatement() ast.Statement {
  // parse according to the type of the current token
  switch p.curToken.Type {
  case token.LET:
    return p.parseLetStatement()
  case token.RETURN:
    return p.parseReturnStatement()
  default:
    return p.parseExpressionStatement()
  }
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
  // construct a let statement node with the current token
  stmt := &ast.LetStatement{Token: p.curToken}

  // expect an identifier
  if !p.expectPeek(token.IDENT) {
    return nil
  }
  // construct an identifier node
  stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

  // expect an equal sign
  if !p.expectPeek(token.ASSIGN) {
    return nil
  }

  // TODO

  // advance until end of statement
  for !p.curTokenIs(token.SEMICOLON) {
    p.nextToken()
  }

  return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
  // construct a return statement node with the current token
  stmt := &ast.ReturnStatement{Token: p.curToken}

  p.nextToken()

  // TODO

  for !p.curTokenIs(token.SEMICOLON) {
    p.nextToken()
  }

  return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	
	// expressions have optional semicolons
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// check if there's a parsing fn associated with current token type in prefix position
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	// if so, call it
	leftExp := prefix()
	return leftExp
}
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:			p.curToken,
		Operator:		p.curToken.Literal,
	}

	// to parse a prefix expression, more than one token needs to be consumed
	// advance tokens and call parseExpression again
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	
	// convert string in current token literal to an int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// return an identifier for current token
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
  return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
  return p.peekToken.Type == t
}

// enforce correctness of token order by checking type of next token
func (p *Parser) expectPeek(t token.TokenType) bool {
  // advance only if peekToken has correct type
  if p.peekTokenIs(t) {
    p.nextToken()
    return true
  } else {
    // add an error every time an expectation about the next token was wrong
    p.peekError(t)
    return false
  }
}

func (p *Parser) Errors() []string {
  return p.errors
}

// add error if type of peekToken doesn't match expectation
func (p *Parser) peekError(t token.TokenType) {
  msg := fmt.Sprintf("expected next token to be %s, got %s", t, p.peekToken.Type)
  p.errors = append(p.errors, msg)
}

// add entries to parser function maps
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
  p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registeInfix(tokenType token.TokenType, fn infixParseFn) {
  p.infixParseFns[tokenType] = fn
}

// define prefix and infix parsing functions
type (
  prefixParseFn func() ast.Expression // both function types return the same type
  infixParseFn  func(ast.Expression) ast.Expression // fn(left_expression) right_expression
)

// operator precedences (increasing order)
const (
	_ int = iota // iota gives these constants incrementing numbers as values
	LOWEST
	EQUALS			// ==
	LESSGREATER	// >, <
	SUM					// +
	PRODUCT			// *
	PREFIX			// -x, !x, etc.
	CALL				// func(X)
)
