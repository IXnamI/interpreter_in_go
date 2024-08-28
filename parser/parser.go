package parser

import (
	"fmt"
	"strconv"

	"github.com/IXnamI/interpreter_in_go/ast"
	"github.com/IXnamI/interpreter_in_go/lexer"
	"github.com/IXnamI/interpreter_in_go/token"
)

type Parser struct {
	l                *lexer.Lexer
	curToken         token.Token
	peekToken        token.Token
	errors           []string
	prefixParseFuncs map[token.TokenType]prefixParseFunc
	infixParseFuncs  map[token.TokenType]infixParseFunc
}

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	p.prefixParseFuncs = make(map[token.TokenType]prefixParseFunc)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.infixParseFuncs = make(map[token.TokenType]infixParseFunc)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		program.Statements = append(program.Statements, statement)
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
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
	letStatement := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	letStatement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return letStatement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStatement := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return returnStatement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expressionStatement := &ast.ExpressionStatement{Token: p.curToken}
	expressionStatement.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return expressionStatement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFuncs[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && (precedence < p.peekPrecedence()) {
		infix := p.infixParseFuncs[p.curToken.Type]
		if infix == nil {
			return leftExpression
		}

		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	intLiteral := &ast.IntegerLiteral{Token: p.curToken}
	intval, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as an int", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}
	intLiteral.Value = intval
	return intLiteral
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(leftExp ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     leftExp,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	p.parseExpression(precedence)

	return expression
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFunc) {
	p.prefixParseFuncs[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFunc) {
	p.infixParseFuncs[tokenType] = fn
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
