package arch

type Opcode uint8

const (
	None Opcode = 0

	ADD  = 1
	ADDI = 2
	SUB  = 3
	SUBI = 4
	AND  = 5
	OR   = 6
	XOR  = 7
	LSL  = 8
	LSR  = 9

	LDREG  = 20
	LDWORD = 21
	LDHWRD = 22
	LDBYTE = 23
	STREG  = 24
	STWORD = 25
	STHWRD = 26
	STBYTE = 27

	MOVZ = 30
	MOVK = 31

	B     = 32
	BREG  = 33
	BLR   = 34
	B_EQ  = 35
	B_NEQ = 36
	B_LT  = 37
	B_LE  = 38
	B_GT  = 39
	B_GE  = 40
	BL    = 41

	PUSH = 42
	POP  = 43

	SYSCALL = 61
	HALT    = 62
	NOOP    = 63
)

type InstructionType uint8

const (
	IType_Invalid InstructionType = iota
	IType_A
	IType_AI
	IType_M
	IType_E
	IType_B
	IType_BI
	IType_R
	IType_O
)

func GetInstructionType(opcode Opcode) InstructionType {
	switch opcode {
	case ADD,
		SUB,
		AND,
		OR,
		XOR:
		return IType_A
	case ADDI,
		SUBI,
		LSL,
		LSR:
		return IType_AI
	case LDREG,
		LDWORD,
		LDHWRD,
		LDBYTE,
		STREG,
		STWORD,
		STHWRD,
		STBYTE:
		return IType_M
	case MOVZ,
		MOVK:
		return IType_E
	case B,
		BL,
		B_EQ,
		B_NEQ,
		B_LT,
		B_LE,
		B_GT,
		B_GE:
		return IType_BI
	case BREG, BLR:
		return IType_B
	case PUSH, POP:
		return IType_R
	case SYSCALL,
		HALT,
		NOOP:
		return IType_O
	default:
		return IType_Invalid
	}
}

func (o Opcode) String() string {
	switch o {
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
	case BLR:
		return "BLR"
	case B_EQ:
		return "B_EQ"
	case B_NEQ:
		return "B_NEQ"
	case B_LT:
		return "B_LT"
	case B_LE:
		return "B_LE"
	case B_GT:
		return "B_GT"
	case B_GE:
		return "B_GE"
	case BL:
		return "BL"
	case PUSH:
		return "PUSH"
	case POP:
		return "POP"
	case SYSCALL:
		return "SYSCALL"
	case HALT:
		return "HALT"
	case NOOP:
		return "NOOP"
	default:
		return "UNKNOWN"
	}
}
