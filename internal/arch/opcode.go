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
	CMP  = 10
	CMPI = 11

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
	B_EQ  = 34
	B_NEQ = 35
	B_LT  = 36
	B_LE  = 37
	B_GT  = 38
	B_GE  = 39
	BL    = 40

	HALT = 62
	NOOP = 63
)
