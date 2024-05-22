package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/IXnamI/interpreter_in_go/lexer"
	"github.com/IXnamI/interpreter_in_go/token"
)

const PROMPT = ">>"

func StartEval(channelIn io.Reader, channelOut io.Writer) {
	scanner := bufio.NewScanner(channelIn)
	for {
		fmt.Fprint(channelOut, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input := scanner.Text()
		l := lexer.New(input)
		for curToken := l.NextToken(); curToken.Type != token.EOF; curToken = l.NextToken() {
			fmt.Fprintf(channelOut, "%+v\n", curToken)
		}
	}
}
