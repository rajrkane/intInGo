// ast/ast.go

package ast

import "intInGo/token"

type Node interface {
  TokenLiteral() string // literal of associated token
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

// root node of every AST
type Program struct {
  Statements []Statement
}

func (p *Program) TokenLiteral() string {
  if len(p.Statements) > 0 {
    return p.Statements[0].TokenLiteral()
  } else {
      return ""
  }
}

// node for let statement (let x = 5)
type LetStatement struct {
  Token token.Token // token.LET
  Name *Identifier  // identifier of the binding (x)
  Value Expression  // expression producing the value (5)
}

func (ls *LetStatement) statementNode() {

}

func (ls *LetStatement) TokenLiteral() string {
  return ls.Token.Literal
}

// identifier (x) of a binding (let x = 5)
type Identifier struct {
  Token token.Token // token.IDENT
  Value string
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) TokenLiteral() string {
  return i.Token.Literal
}
