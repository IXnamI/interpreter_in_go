package lexer

import (
	"strings"

	"github.com/IXnamI/interpreter_in_go/token"
	"github.com/IXnamI/interpreter_in_go/utils"
)

type Lexer struct {
	input   string
	curPos  int
	nextPos int
	curChar byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.curChar = 0
	} else {
		l.curChar = l.input[l.nextPos]
	}
	l.curPos = l.nextPos
	l.nextPos++
}

func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPos]
	}
}

func (l *Lexer) NextToken() token.Token {
	var curToken token.Token
	l.skipWhiteSpaces()
	stringCurChar := string(l.curChar)
	switch l.curChar {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			curToken = token.CreateNewToken(token.EQ, stringCurChar+string(l.curChar))
		} else {
			curToken = token.CreateNewToken(token.ASSIGN, stringCurChar)
		}
	case '+':
		curToken = token.CreateNewToken(token.PLUS, stringCurChar)
	case '-':
		curToken = token.CreateNewToken(token.MINUS, stringCurChar)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			curToken = token.CreateNewToken(token.NOT_EQ, stringCurChar+string(l.curChar))
		} else {
			curToken = token.CreateNewToken(token.BANG, stringCurChar)
		}
	case '*':
		curToken = token.CreateNewToken(token.ASTERISK, stringCurChar)
	case '/':
		curToken = token.CreateNewToken(token.SLASH, stringCurChar)
	case '<':
		curToken = token.CreateNewToken(token.LT, stringCurChar)
	case '>':
		curToken = token.CreateNewToken(token.GT, stringCurChar)
	case ';':
		curToken = token.CreateNewToken(token.SEMICOLON, stringCurChar)
	case '(':
		curToken = token.CreateNewToken(token.LPAREN, stringCurChar)
	case ')':
		curToken = token.CreateNewToken(token.RPAREN, stringCurChar)
	case ',':
		curToken = token.CreateNewToken(token.COMMA, stringCurChar)
	case '{':
		curToken = token.CreateNewToken(token.LBRACE, stringCurChar)
	case '}':
		curToken = token.CreateNewToken(token.RBRACE, stringCurChar)
	case 0:
		curToken.Literal = ""
		curToken.Type = token.EOF
	default:
		if utils.IsLetter(l.curChar) {
			curToken.Literal = l.readGrouped(utils.IsLetter)
			curToken.Type = token.MatchIdentifierTypeFromLiteral(curToken.Literal)
			return curToken
		} else if utils.IsDigit(l.curChar) {
			curToken.Literal = l.readGrouped(utils.IsNumberFormat)
			curToken.Type = validateNumberType(curToken.Literal)
			if curToken.Type == token.BINARY || curToken.Type == token.OCTAL {
				curToken.Literal = curToken.Literal[2:]
			}
			return curToken
		} else {
			curToken = token.CreateNewToken(token.ILLEGAL, stringCurChar)
		}
	}
	l.readChar()
	return curToken
}

func (l *Lexer) readGrouped(checkingFunc func(chr byte) bool) string {
	firstLetterPos := l.curPos
	for checkingFunc(l.curChar) {
		l.readChar()
	}
	return l.input[firstLetterPos:l.curPos]
}

func (l *Lexer) skipWhiteSpaces() {
	for l.curChar == ' ' || l.curChar == '\t' || l.curChar == '\n' || l.curChar == '\r' {
		l.readChar()
	}
}

func validateNumberType(literal string) token.TokenType {
	switch {
	case strings.HasPrefix(literal, "0x"):
		if utils.ValidateOctalNotation(literal) {
			return token.OCTAL
		}
		return token.ILLEGAL
	case strings.HasPrefix(literal, "0b"):
		if utils.ValidateBinaryNotation(literal) {
			return token.BINARY
		}
		return token.ILLEGAL
	case utils.IsFloat(literal):
		return token.FLOAT
	default:
		return token.INT
	}
}
