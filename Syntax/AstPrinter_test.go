package Syntax

import (
	"fmt"
	"github.com/trueabc/lox/Token"
	"testing"
)

func TestAstPrinter_Print(t *testing.T) {
	//(* (- 123) (group 45.67))
	//Expr expression = new Expr.Binary(
	//	new Expr.Unary(
	//	new Token(TokenType.MINUS, "-", null, 1),
	//	new Expr.Literal(123)),
	//new Token(TokenType.STAR, "*", null, 1),
	//	new Expr.Grouping(
	//	new Expr.Literal(45.67)));
	t1 := Token.NewToken(Token.MINUS, "-", nil, 1)
	e1 := &Literal{123}
	u1 := &Unary{t1, e1}
	t2 := Token.NewToken(Token.STAR, "*", nil, 1)
	e2 := &Grouping{&Literal{45.67}}
	e3 := &Binary{u1, t2, e2}

	res := AstPrinter{}.Print(e3)
	fmt.Println(res)
}
