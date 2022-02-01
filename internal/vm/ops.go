package vm

import "math/bits"

func (v *VirtualMachine) add(a, b uint64) uint64 {
	res, carry := bits.Add64(a, b, 0)
	v.setFlags(res, carry)
	return res
}

func (v *VirtualMachine) sub(a, b uint64) uint64 {
	res, _ := bits.Sub64(a, b, 0)
	v.setFlags(res, 0)
	return res
}

func (v *VirtualMachine) and(a, b uint64) uint64 {
	res := a & b
	v.setFlags(res, 0)
	return res
}

func (v *VirtualMachine) or(a, b uint64) uint64 {
	res := a | b
	v.setFlags(res, 0)
	return res
}

func (v *VirtualMachine) xor(a, b uint64) uint64 {
	res := a ^ b
	v.setFlags(res, 0)
	return res
}

func (v *VirtualMachine) lsl(a, b uint64) uint64 {
	res := a << b
	v.setFlags(res, 0)
	return res
}

func (v *VirtualMachine) lsr(a, b uint64) uint64 {
	res := a >> b
	v.setFlags(res, 0)
	return res
}
