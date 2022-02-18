package vm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"math"
)

func (v *VirtualMachine) executeATypeInstruction(instruction arch.ATypeInstruction) {
	aVal := v.registers.Get(instruction.RegA)
	bVal := v.registers.Get(instruction.RegB)

	var res uint64
	switch instruction.Opcode {
	case arch.ADD:
		res = v.add(aVal, bVal)
	case arch.SUB:
		res = v.sub(aVal, bVal)
	case arch.AND:
		res = v.and(aVal, bVal)
	case arch.OR:
		res = v.or(aVal, bVal)
	case arch.XOR:
		res = v.xor(aVal, bVal)
	default:
		panic("invalid A-Type opcode")
	}

	v.registers.Set(instruction.RegDest, res)
}

func (v *VirtualMachine) executeATypeImmInstruction(instruction arch.ATypeImmInstruction) {
	aVal := v.registers.Get(instruction.RegA)
	bVal := uint64(instruction.Immediate)

	var res uint64
	switch instruction.Opcode {
	case arch.ADDI:
		res = v.add(aVal, bVal)
	case arch.SUBI:
		res = v.sub(aVal, bVal)
	case arch.LSL:
		res = v.lsl(aVal, bVal)
	case arch.LSR:
		res = v.lsr(aVal, bVal)
	default:
		panic("invalid AImm-Type opcode")
	}

	v.registers.Set(instruction.RegDest, res)
}

func (v *VirtualMachine) executeMTypeInstruction(instruction arch.MTypeInstruction) {
	baseReg := v.registers.Get(instruction.RegB)
	offset := uint64(instruction.Immediate)

	targetAddress := baseReg + offset
	if targetAddress > math.MaxInt32 {
		panic(fmt.Sprintf("invalid computed memory address %d", targetAddress))
	}

	switch instruction.Opcode {
	case arch.LDREG:
		v.registers.Set(instruction.RegA, v.memory.Read(uint32(targetAddress), 64))
	case arch.LDWORD:
		v.registers.Set(instruction.RegA, v.memory.Read(uint32(targetAddress), 32))
	case arch.LDHWRD:
		v.registers.Set(instruction.RegA, v.memory.Read(uint32(targetAddress), 16))
	case arch.LDBYTE:
		v.registers.Set(instruction.RegA, v.memory.Read(uint32(targetAddress), 8))
	case arch.STREG:
		v.memory.Write(uint32(targetAddress), 64, v.registers.Get(instruction.RegA))
	case arch.STWORD:
		v.memory.Write(uint32(targetAddress), 32, v.registers.Get(instruction.RegA))
	case arch.STHWRD:
		v.memory.Write(uint32(targetAddress), 16, v.registers.Get(instruction.RegA))
	case arch.STBYTE:
		v.memory.Write(uint32(targetAddress), 8, v.registers.Get(instruction.RegA))
	default:
		panic("invalid M-Type opcode")
	}
}

func (v *VirtualMachine) executeETypeInstruction(instruction arch.ETypeInstruction) {
	switch instruction.Opcode {
	case arch.MOVZ:
		v.registers.Set(instruction.RegDest, uint64(instruction.Immediate))
	case arch.MOVK:
		ref := v.registers.Ref(instruction.RegDest)
		*ref = *ref & (0xFFFFFFFFFFFF0000) // clear lower 16 bits
		*ref |= uint64(instruction.Immediate)
	}
}

func (v *VirtualMachine) executeBTypeInstruction(instruction arch.BTypeInstruction) {
	destAddress := v.registers.Get(instruction.RegA)
	switch instruction.Opcode {
	case arch.BREG:
		v.programCounter = uint32(destAddress) - 4
	case arch.BLR:
		nextPC := v.programCounter + 4
		v.registers.Set(arch.ReturnRegister, uint64(nextPC))
		v.programCounter = uint32(destAddress) - 4
	default:
		panic("invalid B-Type opcode")
	}
}

func (v *VirtualMachine) executeBTypeImmInstruction(instruction arch.BTypeImmInstruction) {
	// todo: Verify behavior
	doBranch := false
	doLink := false

	switch instruction.Opcode {
	case arch.B:
		doBranch = true
	case arch.BL:
		doBranch = true
		doLink = true
	case arch.B_EQ:
		doBranch = v.flags.Zero
	case arch.B_NEQ:
		doBranch = !v.flags.Zero
	case arch.B_LT:
		doBranch = v.flags.Negative != v.flags.Carry
	case arch.B_LE:
		doBranch = !(v.flags.Zero && v.flags.Negative == v.flags.Carry)
	case arch.B_GT:
		doBranch = v.flags.Zero && v.flags.Negative == v.flags.Carry
	case arch.B_GE:
		doBranch = v.flags.Negative == v.flags.Carry
	default:
		panic("invalid BImm-Type opcode")
	}

	if doLink {
		nextPC := v.programCounter + 4
		v.registers.Set(arch.ReturnRegister, uint64(nextPC))
	}

	if doBranch {
		v.programCounter += (uint32(instruction.Offset) - 1) * 4
	}
}

func (v *VirtualMachine) executeRTypeInstruction(instruction arch.RTypeInstruction) {
	switch instruction.Opcode {
	case arch.PUSH: // store register then decrement stack pointer
		val := v.registers.Get(instruction.RegA)
		sp := v.registers.Get(arch.StackRegister)
		sp -= 8
		v.registers.Set(arch.StackRegister, sp) // decrement by 8 bytes = 1 register
		v.memory.Write(uint32(sp), 64, val)
	case arch.POP:
		sp := v.registers.Get(arch.StackRegister)
		v.registers.Set(arch.StackRegister, sp+8) // increment by 8 bytes = 1 register
		v.registers.Set(instruction.RegA, v.memory.Read(uint32(sp), 64))
	default:
		panic("invalid R-Type opcode")
	}
}

func (v *VirtualMachine) executeOTypeInstruction(instruction arch.OTypeInstruction) {
	switch instruction.Opcode {
	case arch.SYSCALL:
		v.executeSyscall()
	case arch.NOOP:
		return
	case arch.HALT:
		v.Halt()
	}
}
