// ast/ast.go

package ast

import (
  "intInGo/token"
  "bytes"
)

type Node interface {
  TokenLiteral()  string // literal of associated token
  String()        string // print node for debugging
}
// some nodes implement statement interface
type Statement interface {
  Node
  statementNode()
}
// other nodes implement expression interface
type Expression interface {
  Node
  expressionNode()
}

// program is root node of every AST
type Program struct {
  Statements []Statement
}
func (p *Program) String() string {
  // write return value of each statement's String() to a buffer
  var out bytes.Buffer
  for _, s := range p.Statements {
    out.WriteString(s.String())
  }
  // return buffer as string
  return out.String()
}
func (p *Program) TokenLiteral() string {
  if len(p.Statements) > 0 {
    return p.Statements[0].TokenLiteral()
  } else {
      return ""
  }
}

// identifier (x) of a binding (let x = 5)
type Identifier struct {
  Token token.Token // token.IDENT
  Value string
}
func (i *Identifier) expressionNode() {/* TODO */}
func (i *Identifier) TokenLiteral() string {return i.Token.Literal}
func (i *Identifier) String() string {return i.Value}

// node for let statement (let x = 5)
type LetStatement struct {
  Token token.Token // token.LET
  Name  *Identifier  // identifier of the binding (x)
  Value Expression  // expression producing the value (5)
}
func (ls *LetStatement) statementNode() {/* TODO */}
func (ls *LetStatement) TokenLiteral() string {return ls.Token.Literal}
// print the let statement nicely
func (ls *LetStatement) String() string {
  var out bytes.Buffer
  out.WriteString(ls.TokenLiteral() + " ")
  out.WriteString(ls.Name.String())
  out.WriteString(" = ")
  if ls.Value != nil {
    out.WriteString(ls.Value.String())
  }
  out.WriteString(";")
  return out.String()
}

// node for return statement
type ReturnStatement struct {
  Token       token.Token // token.RETURN
  ReturnValue Expression  // return value
}
func (rs *ReturnStatement) statementNode() {/* TODO */}
func (rs *ReturnStatement) TokenLiteral() string {return rs.Token.Literal}
// print the return statement nicely
func (rs *ReturnStatement) String() string {
  var out bytes.Buffer
  out.WriteString(rs.TokenLiteral() + " ")
  if rs.ReturnValue != nil {
    out.WriteString(rs.ReturnValue.String())
  }
  out.WriteString(";")
  return out.String()
}

// node for expression statement (x + 5;)
type ExpressionStatement struct {
  Token       token.Token // first token of expression
  Expression  Expression
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {return es.Token.Literal}
// print expressiont statement nicely
func (es *ExpressionStatement) String() string {
  if es.Expression != nil {
    return es.Expression.String()
  }
  return ""
}

type IntegerLiteral struct {
	Token		token.Token
	Value		int64 // actual value (not string) of the integer literal
}
func (il *IntegerLiteral) expressionNode() {/* TODO */}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type PrefixExpression struct {
	Token			token.Token // the prefix token (! or -)
	Operator	string
	Right			Expression
}
func (pe *PrefixExpression) expressionNode() {/*TODO*/}
func (pe *PrefixExpression) TokenLiteral() string {return pe.Token.Literal}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token			token.Token // the infix token (e.g. +)
	Left			Expression
	Operator	string
	Right			Expression
}
func (ie *InfixExpression) expressionNode() {/*TODO*/}
func (ie *InfixExpression) TokenLiteral() string {return ie.Token.Literal}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

