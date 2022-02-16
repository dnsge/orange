package lexer

import (
	"fmt"
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"
	"log"
)

func tokenOfKind(kind TokenKind) lexmachine.Action {
	return func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return &Token{
			Kind:   kind,
			Value:  string(match.Bytes),
			Row:    match.StartLine,
			Column: match.StartColumn,
		}, nil
	}
}

func tokenOfKindSliced(kind TokenKind, startOffset, endOffset int) lexmachine.Action {
	return func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		if len(match.Bytes) < startOffset+endOffset {
			return nil, fmt.Errorf("expected string with len >= %d but got %d", startOffset+endOffset, len(match.Bytes))
		}

		return &Token{
			Kind:   kind,
			Value:  string(match.Bytes[startOffset : len(match.Bytes)-endOffset]),
			Row:    match.StartLine,
			Column: match.StartColumn,
		}, nil
	}
}

var sharedLexer *lexmachine.Lexer = nil

func init() {
	// Initialize our assembly lexer with patterns, storing it in sharedLexer
	lexer := lexmachine.NewLexer()

	// Add generated patterns
	addLexerPatterns(lexer)

	lexer.Add([]byte(`[ \t\r]`), func(*lexmachine.Scanner, *machines.Match) (interface{}, error) {
		// skip
		return nil, nil
	})

	err := lexer.CompileDFA()
	if err != nil {
		log.Fatalf("compile lexer dfa: %v\n", err)
	}

	sharedLexer = lexer
}
