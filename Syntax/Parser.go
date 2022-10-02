package Syntax

import (
	"fmt"
	"github.com/trueabc/lox/Errors"
	"github.com/trueabc/lox/Logger"
	"github.com/trueabc/lox/Token"
)

type Parser struct {
	tokens  []*Token.Token
	current int // 下一个需要去消费的token
}

// 还需要检查错误
// 标记恢复点, 在出现error语法之后找到一个恢复点可以继续语法解析
// todo 当前还没有statement的概念, 后面有panic mode, 恢复点以statement分隔

func NewParser(tokens []*Token.Token) *Parser {
	p := &Parser{tokens: tokens, current: 0}
	return p
}

func (p *Parser) Parse() []Stmt {
	defer func() {
		if r := recover(); r != nil {
			Logger.Errorf("%v", r)
		}
	}()
	stmts := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts
}

func (p *Parser) isAtEnd() bool {
	return Token.EOF == p.peek().TType
}

func (p *Parser) peek() *Token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(tokenType ...Token.TokenType) bool {
	for _, item := range tokenType {
		if p.check(item) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t Token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return t == p.peek().TType
}

func (p *Parser) advance() *Token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) assignment() Expr {
	//var a = "before";
	// a = "value";
	expr := p.equality()      // 左侧表达式匹配之后
	if p.match(Token.EQUAL) { // 如果下一个是equal, 说明左侧不应该求值, 而作为token表示符号
		equals := p.previous()
		value := p.assignment()
		if v, ok := expr.(*VariableExpr); ok {
			name := v.name
			return &AssignmentExpr{name: name, value: value}
		}
		content := fmt.Sprintf("Invalid assignment target line: %d.", equals.Line)
		panic(NewParseError(content))
	}
	return expr
}

func (p *Parser) statement() Stmt {
	if p.match(Token.PRINT) {
		return p.printStatement()
	}
	if p.match(Token.LEFT_BRACE) {
		return &BlockStmt{
			statements: p.block(),
		}
	}
	return p.expressionStatement()
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *RuntimeError:
				err := r.(*RuntimeError)
				Errors.LoxRuntimeError(err.Token, err.Content)
			}
			// todo synchronized to new line 避免后面的解析失败
			Logger.Errorf("%v", r)
		}
	}()

	if p.match(Token.VAR) {
		return p.varDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(Token.IDENTIFIER, "Expected variable name here.")
	var initializer Expr
	if p.match(Token.EQUAL) {
		initializer = p.expression()
	}
	p.consume(Token.SEMICOLON, "Expect ';' after variable declaration.")
	return &VariableStmt{name: name, initializer: initializer}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(Token.SEMICOLON, "Expect ';' after value.")
	return &PrintStmt{value}
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)
	for !p.check(Token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(Token.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) expressionStatement() Stmt {
	value := p.expression()
	p.consume(Token.SEMICOLON, "Expect ';' after value.")
	return &ExpressionStmt{value}
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(Token.BANG_EQUAL, Token.EQUAL_EQUAL) {
		op := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{left: expr, operator: op, right: right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(Token.GREATER, Token.GREATER_EQUAL, Token.LESS, Token.LESS_EQUAL) {
		op := p.previous()
		right := p.term()
		expr = &BinaryExpr{left: expr, operator: op, right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(Token.MINUS, Token.PLUS) {
		op := p.previous()
		right := p.factor()
		expr = &BinaryExpr{left: expr, operator: op, right: right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(Token.SLASH, Token.STAR) {
		op := p.previous()
		right := p.unary()
		expr = &BinaryExpr{expr, op, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(Token.BANG, Token.MINUS) {
		op := p.previous()
		right := p.unary()
		return &UnaryExpr{op, right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(Token.FALSE) {
		return &LiteralExpr{false}
	}
	if p.match(Token.TRUE) {
		return &LiteralExpr{true}
	}
	if p.match(Token.NIL) {
		return &LiteralExpr{nil}
	}

	if p.match(Token.NUMBER, Token.STRING) {
		return &LiteralExpr{p.previous().Literal}
	}
	if p.match(Token.LEFT_PAREN) {
		expr := p.expression()

		p.consume(Token.RIGHT_PAREN, "Expect ')' after Expression")

		return &GroupingExpr{expression: expr}
	}
	if p.match(Token.IDENTIFIER) {
		return &VariableExpr{name: p.previous()}
	}
	// 最终匹配到terminal符号, 失败说明当前不是合法的表达式
	panic(p.error(p.peek(), "Expect Expression"))
}

func (p *Parser) consume(tokenType Token.TokenType, mess string) *Token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(p.error(p.peek(), mess))
}

func (p *Parser) error(token *Token.Token, mess string) interface{} {
	data := Errors.LoxError(token, mess)
	return NewParseError(data)
}

type ParseError struct {
	content string
}

func (pe ParseError) Error() string {
	return pe.content
}

func NewParseError(content string) ParseError {
	return ParseError{content: content}
}
