package asm

import (
	"github.com/dnsge/orange/internal/arch"
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

	opKindExpectationMap = map[arch.InstructionType]*lexer.Expectation{
		arch.IType_A:    aType_expectation,
		arch.IType_AI:   aiType_expectation,
		arch.IType_M:    mType_expectation,
		arch.IType_E:    eType_expectation,
		arch.IType_B:    bType_expectation,
		arch.IType_BI:   biType_expectation,
		arch.IType_O:    oType_expectation,
		arch.IType_CMP:  cmpType_expectation,
		arch.IType_CMPI: cmpiType_expectation,
	}
)
