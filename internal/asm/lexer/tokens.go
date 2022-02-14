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
	LBRACKET
	RBRACKET
	COMMENT
	LINE_END

	STRING

	_directiveStart
	LABEL_DECLARATION
	FILL_STATEMENT
	STRING_STATEMENT
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
	ADR

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

// slice of literal token patterns for lexer
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
	ADR:    []byte(`ADR`),

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

// Token describes a lexeme within an input
type Token struct {
	Kind   TokenKind
	Value  string
	Row    int
	Column int
}

// IsImmediate returns whether the TokenKind represents an immediate value
func IsImmediate(kind TokenKind) bool {
	return kind == BASE_10_IMM || kind == BASE_16_IMM || kind == BASE_8_IMM
}

// IsDirective returns whether the TokenKind represents a directive
func IsDirective(kind TokenKind) bool {
	return kind > _directiveStart && kind < _directiveEnd
}

// IsDataDirective returns whether the TokenKind represents a directive that
// will appear as data in the final assembled binary
func IsDataDirective(kind TokenKind) bool {
	return kind == FILL_STATEMENT || kind == STRING_STATEMENT
}

// IsOp returns whether the TokenKind represents an opcode
func IsOp(kind TokenKind) bool {
	return kind > _opStart && kind < _opEnd
}

// GetTokenOpOpcode returns the arch.Opcode for the given TokenKind
func GetTokenOpOpcode(opKind TokenKind) arch.Opcode {
	switch opKind {
	case ADD:
		return arch.ADD
	case ADDI:
		return arch.ADDI
	case SUB:
		return arch.SUB
	case SUBI:
		return arch.SUBI
	case AND:
		return arch.AND
	case OR:
		return arch.OR
	case XOR:
		return arch.XOR
	case LSL:
		return arch.LSL
	case LSR:
		return arch.LSR
	case CMP:
		return arch.CMP
	case CMPI:
		return arch.CMPI
	case LDREG:
		return arch.LDREG
	case LDWORD:
		return arch.LDWORD
	case LDHWRD:
		return arch.LDHWRD
	case LDBYTE:
		return arch.LDBYTE
	case STREG:
		return arch.STREG
	case STWORD:
		return arch.STWORD
	case STHWRD:
		return arch.STHWRD
	case STBYTE:
		return arch.STBYTE
	case MOVZ:
		return arch.MOVZ
	case MOVK:
		return arch.MOVK
	case B:
		return arch.B
	case BREG:
		return arch.BREG
	case B_EQ:
		return arch.B_EQ
	case B_NEQ:
		return arch.B_NEQ
	case B_LT:
		return arch.B_LT
	case B_LE:
		return arch.B_LE
	case B_GT:
		return arch.B_GT
	case B_GE:
		return arch.B_GE
	case BL:
		return arch.BL
	case HALT:
		return arch.HALT
	case NOOP:
		return arch.NOOP
	default:
		panic("lexer.GetTokenOpOpcode: invalid opcode type")
	}
}
