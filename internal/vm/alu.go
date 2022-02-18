package vm

import "math/bits"

type aluFlags struct {
	Negative bool
	Zero     bool
	Carry    bool
}

type ALU struct {
	flags aluFlags
}

func newALU() *ALU {
	return &ALU{
		flags: aluFlags{
			Negative: false,
			Zero:     false,
			Carry:    false,
		},
	}
}

func (alu *ALU) ADD(a, b uint64) uint64 {
	res, carry := bits.Add64(a, b, 0)
	alu.setFlags(res, carry)
	return res
}

func (alu *ALU) SUB(a, b uint64) uint64 {
	res, _ := bits.Sub64(a, b, 0)
	alu.setFlags(res, 0)
	return res
}

func (alu *ALU) AND(a, b uint64) uint64 {
	res := a & b
	alu.setFlags(res, 0)
	return res
}

func (alu *ALU) OR(a, b uint64) uint64 {
	res := a | b
	alu.setFlags(res, 0)
	return res
}

func (alu *ALU) XOR(a, b uint64) uint64 {
	res := a ^ b
	alu.setFlags(res, 0)
	return res
}

func (alu *ALU) LSL(a, b uint64) uint64 {
	res := a << b
	alu.setFlags(res, 0)
	return res
}

func (alu *ALU) LSR(a, b uint64) uint64 {
	res := a >> b
	alu.setFlags(res, 0)
	return res
}

func (alu *ALU) setFlags(res, carry uint64) {
	alu.flags.Zero = res == 0
	alu.flags.Negative = res&(0b1000<<60) > 0 // check last bit for signed-ness
	alu.flags.Carry = carry == 1
}

func (alu *ALU) Zero() bool {
	return alu.flags.Zero
}

func (alu *ALU) Negative() bool {
	return alu.flags.Negative
}

func (alu *ALU) Carry() bool {
	return alu.flags.Carry
}
