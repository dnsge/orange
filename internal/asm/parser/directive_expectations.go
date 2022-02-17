package parser

import (
	"github.com/dnsge/orange/internal/asm/lexer"
)

var (
	// .section [identifier]
	sectionDeclaration_expectation = lexer.NewExpectation(
		".section identifier",
		lexer.Expect(lexer.IDENTIFIER),
	)
	// [$label]:\n
	labelDeclaration_expectation = lexer.NewExpectation("$label:")
	// .fill [imm]
	fillStatement_expectation = lexer.NewExpectation(
		".fill #imm",
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
		lexer.Expect(lexer.LINE_END),
	)
	// .string "my string"
	stringStatement_expectation = lexer.NewExpectation(
		".string \"My string\"",
		lexer.Expect(lexer.STRING),
		lexer.Expect(lexer.LINE_END),
	)

	directiveKindExpectationMap = map[lexer.TokenKind]*lexer.Expectation{
		lexer.SECTION:           sectionDeclaration_expectation,
		lexer.LABEL_DECLARATION: labelDeclaration_expectation,
		lexer.FILL_STATEMENT:    fillStatement_expectation,
		lexer.STRING_STATEMENT:  stringStatement_expectation,
	}
)
