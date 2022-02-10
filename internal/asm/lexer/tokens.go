package lexer

import (
	"github.com/dnsge/orange/internal/arch"
)

type TokenKind int

const (
	REGISTER TokenKind = iota
	BASE_8_IMM
	BASE_10_IMM
	BASE_16_IMM
	LABEL

	COMMA
	COMMENT
	LINE_END

	_directiveStart
	LABEL_DECLARATION
	_directiveEnd

	_opStart
	ADD
	ADDI
	SUB
	SUBI
	AND
	OR
	XOR
	LSL
	LSR
	CMP
	CMPI

	LDREG
	LDWORD
	LDHWRD
	LDBYTE
	STREG
	STWORD
	STHWRD
	STBYTE

	MOVZ
	MOVK

	B
	BREG
	B_EQ
	B_NEQ
	B_LT
	B_LE
	B_GT
	B_GE
	BL

	HALT
	NOOP
	_opEnd
)

var opTokenPatterns = map[TokenKind][]byte{
	ADD:  []byte(`ADD`),
	ADDI: []byte(`ADDI`),
	SUB:  []byte(`SUB`),
	SUBI: []byte(`SUBI`),
	AND:  []byte(`AND`),
	OR:   []byte(`OR`),
	XOR:  []byte(`XOR`),
	LSL:  []byte(`LSL`),
	LSR:  []byte(`LSR`),
	CMP:  []byte(`CMP`),
	CMPI: []byte(`CMPI`),

	LDREG:  []byte(`LDREG`),
	LDWORD: []byte(`LDWORD`),
	LDHWRD: []byte(`LDHWRD`),
	LDBYTE: []byte(`LDBYTE`),
	STREG:  []byte(`STREG`),
	STWORD: []byte(`STWORD`),
	STHWRD: []byte(`STHWRD`),
	STBYTE: []byte(`STBYTE`),

	MOVZ: []byte(`MOVZ`),
	MOVK: []byte(`MOVK`),

	B:     []byte(`B`),
	BREG:  []byte(`BREG`),
	B_EQ:  []byte(`B\.EQ`),
	B_NEQ: []byte(`B\.NEQ`),
	B_LT:  []byte(`B\.LT`),
	B_LE:  []byte(`B\.LE`),
	B_GT:  []byte(`B\.GT`),
	B_GE:  []byte(`B\.GE`),
	BL:    []byte(`BL`),

	HALT: []byte(`HALT`),
	NOOP: []byte(`NOOP`),
}

type Token struct {
	Kind   TokenKind
	Value  string
	Row    int
	Column int
}

func IsImmediate(kind TokenKind) bool {
	return kind == BASE_10_IMM || kind == BASE_16_IMM || kind == BASE_8_IMM
}

func IsDirective(kind TokenKind) bool {
	return kind > _directiveStart && kind < _directiveEnd
}

func IsOp(kind TokenKind) bool {
	return kind > _opStart && kind < _opEnd
}

func GetTokenOpInstructionType(opKind TokenKind) arch.InstructionType {
	switch opKind {
	case ADD,
		SUB,
		AND,
		OR,
		XOR:
		return arch.IType_A
	case ADDI,
		SUBI,
		LSL,
		LSR:
		return arch.IType_AI
	case LDREG,
		LDWORD,
		LDHWRD,
		LDBYTE,
		STREG,
		STWORD,
		STHWRD,
		STBYTE:
		return arch.IType_M
	case MOVZ,
		MOVK:
		return arch.IType_E
	case B,
		B_EQ,
		B_NEQ,
		B_LT,
		B_LE,
		B_GT,
		B_GE:
		return arch.IType_BI
	case BREG:
		return arch.IType_B
	case HALT,
		NOOP:
		return arch.IType_O
	case CMP:
		return arch.IType_CMP // pseudo-type
	case CMPI:
		return arch.IType_CMPI // pseudo-type

	default:
		return arch.IType_Invalid
	}
}
