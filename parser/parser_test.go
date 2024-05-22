package parser

import (
	"testing"

	"github.com/IXnamI/interpreter_in_go/ast"
	"github.com/IXnamI/interpreter_in_go/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Number of parsed statements does not contain 3 statements, got = %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, testCase := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, testCase.expectedIdentifer) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentiferExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not the right amount of statements, got = %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Parsed statement is not an expression statement, got = %T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("Expression not of type *ast.Identifier, got = %T", stmt.Expression)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("Identifier's token literal not %s, got = %s", "foobar", ident.TokenLiteral())
	}
}

func testLetStatement(t *testing.T, s ast.Statement, expectedValue string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("Statement's token literal not 'let', got = %q", s.TokenLiteral())
		return false
	}
	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement's type not *ast.LetStatement, got = %T", s)
		return false
	}

	if letStatement.Name.Value != expectedValue {
		t.Errorf("Statement's value not %s, got = %s", expectedValue, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != expectedValue {
		t.Errorf("Statement's token literal not '%s', got = %s", expectedValue, letStatement.Name.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return

	}
	t.Errorf("Parser return %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}
