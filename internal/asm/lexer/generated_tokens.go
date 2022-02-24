// Generated token definitions
//
// Generated at 2022-02-24T14:30:33-05:00

package lexer

import (
	"github.com/dnsge/orange/internal/arch"
	"github.com/timtadh/lexmachine"
)

type TokenKind int

const (
	_invalid TokenKind = iota
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
	MOV
	MOVZ
	MOVK
	B
	BREG
	BLR
	B_EQ
	B_NEQ
	B_LT
	B_LE
	B_GT
	B_GE
	BL
	PUSH
	POP
	SYSCALL
	HALT
	NOOP
	_opEnd
	_identifierStart
	REGISTER
	IDENTIFIER
	_identifierEnd
	_immStart
	BASE_8_IMM
	BASE_10_IMM
	BASE_16_IMM
	_immEnd
	_directiveStart
	LABEL_DECLARATION
	SECTION
	FILL_STATEMENT
	STRING_STATEMENT
	ADDRESS_OF
	_directiveEnd
	STRING
	LABEL
	COMMA
	LBRACKET
	RBRACKET
	COMMENT
	LINE_END
)

// IsTokenIdentifier returns whether the token is in the Identifier category
func IsTokenIdentifier(kind TokenKind) bool {
	return kind > _identifierStart && kind < _identifierEnd
}

// IsTokenImm returns whether the token is in the Imm category
func IsTokenImm(kind TokenKind) bool {
	return kind > _immStart && kind < _immEnd
}

// IsTokenDirective returns whether the token is in the Directive category
func IsTokenDirective(kind TokenKind) bool {
	return kind > _directiveStart && kind < _directiveEnd
}

// IsTokenOp returns whether the token is in the Op category
func IsTokenOp(kind TokenKind) bool {
	return kind > _opStart && kind < _opEnd
}

// addLexerPatterns initializes the lexer with patterns for instruction parsing
func addLexerPatterns(lexer *lexmachine.Lexer) {
	// REGISTER
	lexer.Add([]byte("r[0-9]"), tokenOfKind(REGISTER))
	// REGISTER
	lexer.Add([]byte("r1[0-5]"), tokenOfKind(REGISTER))
	// REGISTER
	lexer.Add([]byte("rsp"), tokenOfKind(REGISTER))
	// BASE_8_IMM
	lexer.Add([]byte("#0o(-?[0-7]+)"), tokenOfKind(BASE_8_IMM))
	// BASE_10_IMM
	lexer.Add([]byte("#(0|-?[1-9][0-9]*)"), tokenOfKind(BASE_10_IMM))
	// BASE_16_IMM
	lexer.Add([]byte("#0x(-?[0-9A-Fa-f]+)"), tokenOfKind(BASE_16_IMM))
	// STRING
	lexer.Add([]byte("\"(\\\\\"|[^\"])*\""), tokenOfString(STRING))
	// STRING
	lexer.Add([]byte("`[^`]*`"), tokenOfString(STRING))
	// LABEL_DECLARATION
	lexer.Add([]byte("\\$[a-zA-Z_][a-zA-Z0-9_.]*:"), tokenOfKindSliced(LABEL_DECLARATION, 1, 1))
	// LABEL
	lexer.Add([]byte("\\$[a-zA-Z_][a-zA-Z0-9_.]*"), tokenOfKindSliced(LABEL, 1, 0))
	// SECTION
	lexer.Add([]byte("\\.section"), tokenOfKind(SECTION))
	// FILL_STATEMENT
	lexer.Add([]byte("\\.fill"), tokenOfKind(FILL_STATEMENT))
	// STRING_STATEMENT
	lexer.Add([]byte("\\.string"), tokenOfKind(STRING_STATEMENT))
	// ADDRESS_OF
	lexer.Add([]byte("\\.addressOf"), tokenOfKind(ADDRESS_OF))
	// COMMA
	lexer.Add([]byte(","), tokenOfKind(COMMA))
	// LBRACKET
	lexer.Add([]byte("\\["), tokenOfKind(LBRACKET))
	// RBRACKET
	lexer.Add([]byte("\\]"), tokenOfKind(RBRACKET))
	// COMMENT
	lexer.Add([]byte(";[^\\n]*"), tokenOfKind(COMMENT))
	// LINE_END
	lexer.Add([]byte("\\n"), tokenOfKind(LINE_END))
	// ADD
	lexer.Add([]byte("ADD"), tokenOfKind(ADD))
	// ADDI
	lexer.Add([]byte("ADDI"), tokenOfKind(ADDI))
	// SUB
	lexer.Add([]byte("SUB"), tokenOfKind(SUB))
	// SUBI
	lexer.Add([]byte("SUBI"), tokenOfKind(SUBI))
	// AND
	lexer.Add([]byte("AND"), tokenOfKind(AND))
	// OR
	lexer.Add([]byte("OR"), tokenOfKind(OR))
	// XOR
	lexer.Add([]byte("XOR"), tokenOfKind(XOR))
	// LSL
	lexer.Add([]byte("LSL"), tokenOfKind(LSL))
	// LSR
	lexer.Add([]byte("LSR"), tokenOfKind(LSR))
	// CMP
	lexer.Add([]byte("CMP"), tokenOfKind(CMP))
	// CMPI
	lexer.Add([]byte("CMPI"), tokenOfKind(CMPI))
	// LDREG
	lexer.Add([]byte("LDREG"), tokenOfKind(LDREG))
	// LDWORD
	lexer.Add([]byte("LDWORD"), tokenOfKind(LDWORD))
	// LDHWRD
	lexer.Add([]byte("LDHWRD"), tokenOfKind(LDHWRD))
	// LDBYTE
	lexer.Add([]byte("LDBYTE"), tokenOfKind(LDBYTE))
	// STREG
	lexer.Add([]byte("STREG"), tokenOfKind(STREG))
	// STWORD
	lexer.Add([]byte("STWORD"), tokenOfKind(STWORD))
	// STHWRD
	lexer.Add([]byte("STHWRD"), tokenOfKind(STHWRD))
	// STBYTE
	lexer.Add([]byte("STBYTE"), tokenOfKind(STBYTE))
	// ADR
	lexer.Add([]byte("ADR"), tokenOfKind(ADR))
	// MOV
	lexer.Add([]byte("MOV"), tokenOfKind(MOV))
	// MOVZ
	lexer.Add([]byte("MOVZ"), tokenOfKind(MOVZ))
	// MOVK
	lexer.Add([]byte("MOVK"), tokenOfKind(MOVK))
	// B
	lexer.Add([]byte("B"), tokenOfKind(B))
	// BREG
	lexer.Add([]byte("BREG"), tokenOfKind(BREG))
	// BLR
	lexer.Add([]byte("BLR"), tokenOfKind(BLR))
	// B.EQ
	lexer.Add([]byte("B\\.EQ"), tokenOfKind(B_EQ))
	// B.NEQ
	lexer.Add([]byte("B\\.NEQ"), tokenOfKind(B_NEQ))
	// B.LT
	lexer.Add([]byte("B\\.LT"), tokenOfKind(B_LT))
	// B.LE
	lexer.Add([]byte("B\\.LE"), tokenOfKind(B_LE))
	// B.GT
	lexer.Add([]byte("B\\.GT"), tokenOfKind(B_GT))
	// B.GE
	lexer.Add([]byte("B\\.GE"), tokenOfKind(B_GE))
	// BL
	lexer.Add([]byte("BL"), tokenOfKind(BL))
	// PUSH
	lexer.Add([]byte("PUSH"), tokenOfKind(PUSH))
	// POP
	lexer.Add([]byte("POP"), tokenOfKind(POP))
	// SYSCALL
	lexer.Add([]byte("SYSCALL"), tokenOfKind(SYSCALL))
	// HALT
	lexer.Add([]byte("HALT"), tokenOfKind(HALT))
	// NOOP
	lexer.Add([]byte("NOOP"), tokenOfKind(NOOP))
	// IDENTIFIER
	lexer.Add([]byte("[a-zA-Z][a-zA-Z0-9]*"), tokenOfKind(IDENTIFIER))
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
	case BLR:
		return arch.BLR
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
	case PUSH:
		return arch.PUSH
	case POP:
		return arch.POP
	case SYSCALL:
		return arch.SYSCALL
	case HALT:
		return arch.HALT
	case NOOP:
		return arch.NOOP
	default:
		panic("lexer.GetTokenOpOpcode: invalid opcode type")
	}
}
