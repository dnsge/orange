package vm

import "github.com/dnsge/orange/internal/arch"

func (v *VirtualMachine) fetchNextInstruction() arch.Instruction {
	i := v.memory.Read(v.programCounter, 32) // read word for instruction
	return arch.Instruction(i)
}

func (v *VirtualMachine) executeInstruction(instruction arch.Instruction) {
	opcode := arch.GetOpcode(instruction)
	iType := arch.GetInstructionType(opcode)

	switch iType {
	case arch.IType_A:
		i := arch.DecodeATypeInstruction(instruction, opcode)
		v.executeATypeInstruction(i)
	case arch.IType_AI:
		i := arch.DecodeATypeImmInstruction(instruction, opcode)
		v.executeATypeImmInstruction(i)
	case arch.IType_M:
		i := arch.DecodeMTypeInstruction(instruction, opcode)
		v.executeMTypeInstruction(i)
	case arch.IType_E:
		i := arch.DecodeETypeInstruction(instruction, opcode)
		v.executeETypeInstruction(i)
	case arch.IType_BI:
		i := arch.DecodeBTypeImmInstruction(instruction, opcode)
		v.executeBTypeImmInstruction(i)
	case arch.IType_B:
		i := arch.DecodeBTypeInstruction(instruction, opcode)
		v.executeBTypeInstruction(i)
	case arch.IType_R:
		i := arch.DecodeRTypeInstruction(instruction, opcode)
		v.executeRTypeInstruction(i)
	case arch.IType_O:
		i := arch.DecodeOTypeInstruction(instruction, opcode)
		v.executeOTypeInstruction(i)
	default:
		panic("invalid instruction type")
	}

	v.programCounter += 4 // advance by word
}
