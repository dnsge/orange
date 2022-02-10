package asm

import (
	"github.com/dnsge/orange/internal/asm/lexer"
)

var (
	// [.label]:\n
	labelDeclaration_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.LINE_END),
	)

	directiveKindExpectationMap = map[lexer.TokenKind]*lexer.Expectation{
		lexer.LABEL_DECLARATION: labelDeclaration_expectation,
	}
)
