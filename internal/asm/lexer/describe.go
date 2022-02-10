package lexer

import "fmt"

// DescribeToken returns a human-readable token description,
// including contextual information if applicable like register number.
func DescribeToken(token *Token) string {
	switch token.Kind {
	case REGISTER:
		return token.Value
	case BASE_8_IMM:
		return token.Value
	case BASE_10_IMM:
		return token.Value
	case BASE_16_IMM:
		return token.Value
	case LABEL:
		return token.Value
	case COMMENT:
		return fmt.Sprintf("comment %q", token.Value)
	case LABEL_DECLARATION:
		return token.Value
	default:
		return DescribeTokenKind(token.Kind)
	}
}

// DescribeTokenKind returns a human-readable TokenKind
func DescribeTokenKind(kind TokenKind) string {
	switch kind {
	case REGISTER:
		return "<register>"
	case BASE_8_IMM:
		return "<base 8 imm>"
	case BASE_10_IMM:
		return "<base 10 imm>"
	case BASE_16_IMM:
		return "<base 16 imm>"
	case LABEL:
		return "<label>"
	case COMMA:
		return "<comma>"
	case COMMENT:
		return "<comment>"
	case LINE_END:
		return "<line end>"
	case LABEL_DECLARATION:
		return "<label declaration>"
	case ADD:
		return "ADD"
	case ADDI:
		return "ADDI"
	case SUB:
		return "SUB"
	case SUBI:
		return "SUBI"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	case LSL:
		return "LSL"
	case LSR:
		return "LSR"
	case CMP:
		return "CMP"
	case CMPI:
		return "CMPI"
	case LDREG:
		return "LDREG"
	case LDWORD:
		return "LDWORD"
	case LDHWRD:
		return "LDHWRD"
	case LDBYTE:
		return "LDBYTE"
	case STREG:
		return "STREG"
	case STWORD:
		return "STWORD"
	case STHWRD:
		return "STHWRD"
	case STBYTE:
		return "STBYTE"
	case MOVZ:
		return "MOVZ"
	case MOVK:
		return "MOVK"
	case B:
		return "B"
	case BREG:
		return "BREG"
	case B_EQ:
		return "B.EQ"
	case B_NEQ:
		return "B.NEQ"
	case B_LT:
		return "B.LT"
	case B_LE:
		return "B.LE"
	case B_GT:
		return "B.GT"
	case B_GE:
		return "B.GE"
	case BL:
		return "BL"
	case HALT:
		return "HALT"
	case NOOP:
		return "NOOP"
	default:
		return "<unknown>"
	}
}
