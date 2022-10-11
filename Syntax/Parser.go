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
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	//var a = "before";
	// a = "value";
	expr := p.or()            // 左侧表达式匹配之后
	if p.match(Token.EQUAL) { // 如果下一个是equal, 说明左侧不应该求值, 而作为token表示符号
		equals := p.previous()
		value := p.assignment()
		if v, ok := expr.(*VariableExpr); ok {
			name := v.name
			return &AssignmentExpr{name: name, value: value}
		}
		if v, ok := expr.(*GetExpr); ok {
			return &SetExpr{v.object, v.name, value}
		}
		content := fmt.Sprintf("Invalid assignment target line: %d.", equals.Line)
		panic(NewParseError(content))
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(Token.OR) {
		operator := p.previous()
		right := p.and()
		expr = &LogicExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(Token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &LogicExpr{expr, operator, right}
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
	if p.match(Token.IF) {
		return p.ifStatement()
	}
	if p.match(Token.WHILE) {
		return p.whileStatement()
	}
	if p.match(Token.FOR) {
		return p.forStatement()
	}
	if p.match(Token.RETURN) {
		return p.returnStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(Token.SEMICOLON) {
		value = p.expression()
	}
	p.consume(Token.SEMICOLON, "Expect ; after value.")

	return &ReturnStmt{keyword: keyword, value: value}
}

func (p *Parser) forStatement() Stmt {
	p.consume(Token.LEFT_PAREN, "Expected left '(' after 'for'")
	var initializer Stmt
	if p.match(Token.SEMICOLON) {
		initializer = nil
	} else if p.match(Token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	var condition Expr
	if !p.check(Token.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(Token.SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(Token.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(Token.RIGHT_PAREN, "Expect ')' after for clauses.")
	body := p.statement()

	// 合成while
	if condition == nil {
		condition = &LiteralExpr{true}
	}
	if increment != nil {
		body = &BlockStmt{[]Stmt{body, &ExpressionStmt{increment}}}
	}

	body = &WhileStmt{body: body, condition: condition}

	if initializer != nil {
		body = &BlockStmt{[]Stmt{initializer, body}}
	}
	// 	类似这样的语法
	// {
	//	var i := 1
	// 	while i < 10 {
	//  	print i;
	// 		i += 1;
	// }
	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(Token.LEFT_PAREN, "Expected left '(' after 'while'.")
	condition := p.expression()
	p.consume(Token.RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return &WhileStmt{condition: condition, body: body}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(Token.LEFT_PAREN, "Expected left '(' after 'if'.")
	condition := p.expression()
	p.consume(Token.RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch Stmt

	if p.match(Token.ELSE) {
		elseBranch = p.statement()
	}

	return &IfStmt{condition, thenBranch, elseBranch}
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *RuntimeError:
				err := r.(*RuntimeError)
				Errors.LoxRuntimeError(err.Token, err.Content)
			}
			for p.peek().Lexeme != "\n" && !p.isAtEnd() {
				p.current += 1
			}
			// todo synchronized to new line 避免后面的解析失败
			Logger.Errorf("%v", r)
		}
	}()

	if p.match(Token.VAR) {
		return p.varDeclaration()
	} else if p.match(Token.FUN) {
		return p.function("function")
	} else if p.match(Token.CLASS) {
		return p.classDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(Token.IDENTIFIER, "Expected class name.")

	var superClass *VariableExpr
	if p.match(Token.LESS) {
		p.consume(Token.IDENTIFIER, "Expect superclass name.")
		superClass = &VariableExpr{p.previous()}
	}

	p.consume(Token.LEFT_BRACE, "Expected '{' before class body.")

	// []functionStmt
	methods := make([]Stmt, 0)
	for !p.check(Token.RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}
	p.consume(Token.RIGHT_BRACE, "Expected '}' after class body.")

	return &ClassStmt{name: name, methods: methods, superClass: superClass}
}

func (p *Parser) function(kind string) Stmt {
	name := p.consume(Token.IDENTIFIER, "Expected "+kind+" name.")
	p.consume(Token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	params := make([]*Token.Token, 0)
	if !p.check(Token.RIGHT_PAREN) {
		params = append(params, p.consume(Token.IDENTIFIER, "Expected parameter name."))
		for p.match(Token.COMMA) {
			if len(params) >= 255 {
				p.error(p.peek(), "Can't have more than 255 parameters.")
			}
			params = append(params, p.consume(Token.IDENTIFIER, "Expected parameter name."))
		}
	}
	p.consume(Token.RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(Token.LEFT_BRACE, "Expect '{' before "+kind+" body.")

	body := p.block()
	return &FunctionStmt{name: name, params: params, body: body}
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

	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(Token.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(Token.DOT) {
			name := p.consume(Token.IDENTIFIER,
				"Expect property name after '.'.")
			expr = &GetExpr{expr, name}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	// 获取内部的参数
	arguments := make([]Expr, 0)
	if !p.check(Token.RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
		for p.match(Token.COMMA) {
			if len(arguments) > 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
		}
	}

	paren := p.consume(Token.RIGHT_PAREN,
		"Expect ')' after arguments.")
	// 括号用于进行定位
	return &CallExpr{callee: callee, arguments: arguments, paren: paren}
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
	if p.match(Token.THIS) {
		return &ThisExpr{p.previous()}
	}
	if p.match(Token.SUPER) {
		keyword := p.previous()
		p.consume(Token.DOT, "Expect '.' after 'super'.")
		method := p.consume(Token.IDENTIFIER,
			"Expect superclass method name.")
		return &SuperExpr{keyword, method}
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
