package parser

import (
	"github.com/dnsge/orange/internal/asm/lexer"
)

var (
	// .section [identifier]
	sectionDeclaration_expectation = NewExpectation(
		".section identifier",
		Expect(lexer.IDENTIFIER),
	)
	// [$label]:\n
	labelDeclaration_expectation = NewExpectation("$label:")
	// .fill [imm]
	fillStatement_expectation = OneOf(
		NewExpectation(
			".fill #imm",
			ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
			Expect(lexer.LINE_END),
		),
		NewExpectation(
			".fill .addressOf $label",
			Expect(lexer.ADDRESS_OF),
			Expect(lexer.LABEL),
		),
	)
	// .string "my string"
	stringStatement_expectation = NewExpectation(
		".string \"My string\"",
		Expect(lexer.STRING),
		Expect(lexer.LINE_END),
	)
	// .addressOf $label
	addressOf_expectation = NewExpectation(
		".addressOf $label",
		Expect(lexer.LABEL),
	)

	directiveKindExpectationMap = map[lexer.TokenKind]Extractable{
		lexer.SECTION:           sectionDeclaration_expectation,
		lexer.LABEL_DECLARATION: labelDeclaration_expectation,
		lexer.FILL_STATEMENT:    fillStatement_expectation,
		lexer.STRING_STATEMENT:  stringStatement_expectation,
		lexer.ADDRESS_OF:        addressOf_expectation,
	}
)
