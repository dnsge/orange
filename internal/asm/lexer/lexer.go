package lexer

import (
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

var sharedLexer *lexmachine.Lexer = nil

func init() {
	lexer := lexmachine.NewLexer()
	lexer.Add([]byte(`r[0-9]`), tokenOfKind(REGISTER))
	lexer.Add([]byte(`r1[0-5]`), tokenOfKind(REGISTER))
	lexer.Add([]byte(`#0o(-?[0-7]+)`), tokenOfKind(BASE_8_IMM))
	lexer.Add([]byte(`#(0|-?[1-9][0-9]*)`), tokenOfKind(BASE_10_IMM))
	lexer.Add([]byte(`#0x(-?[0-9A-Fa-f]+)`), tokenOfKind(BASE_16_IMM))

	lexer.Add([]byte(`\.[a-zA-Z][a-zA-Z0-9]*:`), func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
		return &Token{
			Kind:   LABEL_DECLARATION,
			Value:  string(match.Bytes[:len(match.Bytes)-1]), // discard matched colon
			Row:    match.StartLine,
			Column: match.StartColumn,
		}, nil
	})
	lexer.Add([]byte(`\.[a-zA-Z][a-zA-Z0-9]*`), tokenOfKind(LABEL))

	lexer.Add([]byte(`,`), tokenOfKind(COMMA))
	lexer.Add([]byte(`;[^\n]*`), tokenOfKind(COMMENT))

	for tokenKind, pattern := range opTokenPatterns {
		lexer.Add(pattern, tokenOfKind(tokenKind))
	}

	lexer.Add([]byte(`\n`), tokenOfKind(LINE_END))
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