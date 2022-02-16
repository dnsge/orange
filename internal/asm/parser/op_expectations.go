package parser

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
)

var (
	// [OPCODE] [DEST], [SRC1], [SRC2]
	aType_expectation = lexer.NewExpectation(
		"OPCODE r3, r1, r2",
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
	)
	// [OPCODE] [DEST], [SRC1], [IMM]
	aiType_expectation = lexer.NewExpectation(
		"OPCODE r2, r1, #imm",
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
	mType_expectation = lexer.OneOf(
		// [OPCODE] [REG1], [[REG2], [IMM]]
		lexer.NewExpectation(
			"OPCODE r2, [r1, #imm]",
			lexer.Expect(lexer.REGISTER),
			lexer.ExpectIgnore(lexer.COMMA),
			lexer.ExpectIgnore(lexer.LBRACKET),
			lexer.Expect(lexer.REGISTER),
			lexer.ExpectIgnore(lexer.COMMA),
			lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
			lexer.ExpectIgnore(lexer.RBRACKET),
		),
		// [OPCODE] [REG1], [[REG2]] (no offset)
		lexer.NewExpectation(
			"OPCODE r2, [r1]",
			lexer.Expect(lexer.REGISTER),
			lexer.ExpectIgnore(lexer.COMMA),
			lexer.ExpectIgnore(lexer.LBRACKET),
			lexer.Expect(lexer.REGISTER),
			lexer.ExpectIgnore(lexer.RBRACKET),
		),
	)
	// [OPCODE] [DEST], [IMM]
	eType_expectation = lexer.NewExpectation(
		"OPCODE r1, #imm",
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
	// [OPCODE] [REG]
	bType_expectation = lexer.NewExpectation(
		"OPCODE r1",
		lexer.Expect(lexer.REGISTER),
	)
	// [OPCODE] [IMM|LABEL]
	biType_expectation = lexer.NewExpectation(
		"OPCODE [#imm|$label]",
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM, lexer.LABEL),
	)
	// [OPCODE]
	oType_expectation = lexer.NewExpectation("OPCODE")

	// ------ Translated Instructions ------

	// [OPCODE] [REG1], [REG2]
	cmp_expectation = lexer.NewExpectation(
		"CMP r1, r2",
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
	)
	// [OPCODE] [REG1], [IMM]
	cmpi_expectation = lexer.NewExpectation(
		"CMP r1, #imm",
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
	// [OPCODE] [REG1], [$label]
	adr_expectation = lexer.NewExpectation(
		"ADR r1, $label",
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.LABEL),
	)
	// [OPCODE] [REG1]
	r_expectation = lexer.NewExpectation(
		"OPCODE r1",
		lexer.Expect(lexer.REGISTER),
	)
)

func getOpcodeStatementExpectation(opKind lexer.TokenKind) (lexer.Extractable, error) {
	switch opKind {
	case lexer.ADD,
		lexer.SUB,
		lexer.AND,
		lexer.OR,
		lexer.XOR:
		return aType_expectation, nil
	case lexer.ADDI,
		lexer.SUBI,
		lexer.LSL,
		lexer.LSR:
		return aiType_expectation, nil
	case lexer.LDREG,
		lexer.LDWORD,
		lexer.LDHWRD,
		lexer.LDBYTE,
		lexer.STREG,
		lexer.STWORD,
		lexer.STHWRD,
		lexer.STBYTE:
		return mType_expectation, nil
	case lexer.MOVZ,
		lexer.MOVK:
		return eType_expectation, nil
	case lexer.BREG, lexer.BL:
		return bType_expectation, nil
	case lexer.B,
		lexer.B_EQ,
		lexer.B_NEQ,
		lexer.B_LT,
		lexer.B_LE,
		lexer.B_GT,
		lexer.B_GE:
		return biType_expectation, nil
	case lexer.PUSH, lexer.POP:
		return r_expectation, nil
	case lexer.SYSCALL,
		lexer.HALT,
		lexer.NOOP:
		return oType_expectation, nil
	case lexer.CMP:
		return cmp_expectation, nil
	case lexer.CMPI:
		return cmpi_expectation, nil
	case lexer.ADR:
		return adr_expectation, nil
	default:
		return nil, fmt.Errorf("getOpcodeStatementExpectation: invalid opcode")
	}
}
