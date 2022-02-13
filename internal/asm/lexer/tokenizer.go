package lexer

import (
	"github.com/timtadh/lexmachine"
	"io"
)

type Tokenizer struct {
	scanner *lexmachine.Scanner
	input   []byte
}

// New returns a Tokenizer over a byte slice
func New(input []byte) (*Tokenizer, error) {
	scanner, err := sharedLexer.Scanner(input)
	if err != nil {
		return nil, err
	}

	return &Tokenizer{
		scanner: scanner,
		input:   input,
	}, nil
}

// Next extracts the next token from the input buffer.
// Returns the token, an error. Returns io.EOF on end.
func (t *Tokenizer) Next() (*Token, error) {
	tok, err, end := t.scanner.Next()

	if err != nil {
		return nil, err
	} else if end {
		return nil, io.EOF
	}

	return tok.(*Token), nil
}
