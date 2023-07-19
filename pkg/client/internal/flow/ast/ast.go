package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/claion-org/claiflow/pkg/client/internal/flow/token"
)

type Node interface {
	TokenLiteral() string
	String(in map[string]interface{}) string
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String(in map[string]interface{}) string {
	if es.Expression != nil {
		return es.Expression.String(in)
	}
	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String(in map[string]interface{}) string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()                         {}
func (il *IntegerLiteral) TokenLiteral() string                    { return il.Token.Literal }
func (il *IntegerLiteral) String(in map[string]interface{}) string { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()                         {}
func (sl *StringLiteral) TokenLiteral() string                    { return sl.Token.Literal }
func (sl *StringLiteral) String(in map[string]interface{}) string { return sl.Token.Literal }

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String(in map[string]interface{}) string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String(in))
	out.WriteString("[")

	if v, ok := in[strings.TrimPrefix(ie.Index.String(in), "$")]; ok {
		out.WriteString(fmt.Sprintf("%v", v))
	} else {
		out.WriteString(ie.Index.String(in))
	}

	out.WriteString("]")

	return out.String()
}

type SelectorExpression struct {
	Token token.Token
	Left  Expression
	Sel   Expression
}

func (oe *SelectorExpression) expressionNode()      {}
func (oe *SelectorExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *SelectorExpression) String(in map[string]interface{}) string {
	var out bytes.Buffer

	out.WriteString(oe.Left.String(in))
	out.WriteString(".")
	out.WriteString(oe.Sel.String(in))

	return out.String()
}
