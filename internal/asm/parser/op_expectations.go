package parser

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
)

var (
	// [OPCODE] [DEST], [SRC1], [SRC2]
	aType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
	)
	// [OPCODE] [DEST], [SRC1], [IMM]
	aiType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
	// [OPCODE] [REG1], [REG2], [IMM]
	mType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
	// [OPCODE] [DEST], [IMM]
	eType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
	// [OPCODE] [REG]
	bType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
	)
	// [OPCODE] [IMM|LABEL]
	biType_expectation = lexer.NewExpectation(
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM, lexer.LABEL),
	)
	// [OPCODE]
	oType_expectation = lexer.NewExpectation()
	// [OPCODE] [REG1], [REG2]
	cmpType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.Expect(lexer.REGISTER),
	)
	// [OPCODE] [REG1], [IMM]
	cmpiType_expectation = lexer.NewExpectation(
		lexer.Expect(lexer.REGISTER),
		lexer.ExpectIgnore(lexer.COMMA),
		lexer.ExpectAny(lexer.BASE_10_IMM, lexer.BASE_16_IMM, lexer.BASE_8_IMM),
	)
)

func getOpcodeStatementExpectation(opKind lexer.TokenKind) (*lexer.Expectation, error) {
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
	case lexer.B,
		lexer.B_EQ,
		lexer.B_NEQ,
		lexer.B_LT,
		lexer.B_LE,
		lexer.B_GT,
		lexer.B_GE:
		return biType_expectation, nil
	case lexer.BREG:
		return bType_expectation, nil
	case lexer.HALT,
		lexer.NOOP:
		return oType_expectation, nil
	case lexer.CMP:
		return cmpType_expectation, nil
	case lexer.CMPI:
		return cmpiType_expectation, nil
	default:
		return nil, fmt.Errorf("getOpcodeStatementExpectation: invalid opcode")
	}
}
