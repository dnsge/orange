package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"strconv"
)

var (
	ErrInvalidArgumentCount = fmt.Errorf("invalid argument count")
)

func assembleATypeInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	var instruction arch.ATypeInstruction
	if opcode == arch.CMP {
		if len(args) != 2 {
			return 0, ErrInvalidArgumentCount
		}

		parsedRegs, err := parseRegisters(args)
		if err != nil {
			return 0, err
		}

		// CMP is a pseudo-instruction for SUB
		instruction = arch.ATypeInstruction{
			Opcode:  arch.SUB,
			RegDest: 0,
			RegA:    parsedRegs[0],
			RegB:    parsedRegs[1],
		}
	} else {
		if len(args) != 3 {
			return 0, ErrInvalidArgumentCount
		}

		parsedRegs, err := parseRegisters(args)
		if err != nil {
			return 0, err
		}

		instruction = arch.ATypeInstruction{
			Opcode:  opcode,
			RegDest: parsedRegs[0],
			RegA:    parsedRegs[1],
			RegB:    parsedRegs[2],
		}
	}

	return arch.EncodeATypeInstruction(instruction), nil
}

func assembleATypeImmInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	var instruction arch.ATypeImmInstruction
	if opcode == arch.CMPI {
		if len(args) != 2 {
			return 0, ErrInvalidArgumentCount
		}

		regA, err := parseRegister(args[0])
		if err != nil {
			return 0, err
		}

		imm, err := parseUnsignedImmediate(args[1])
		if err != nil {
			return 0, err
		}

		instruction = arch.ATypeImmInstruction{
			Opcode:    arch.SUB,
			RegDest:   0,
			RegA:      regA,
			Immediate: imm,
		}
	} else {
		if len(args) != 3 {
			return 0, ErrInvalidArgumentCount
		}

		parsedRegs, err := parseRegisters(args[0:2])
		if err != nil {
			return 0, err
		}

		imm, err := parseUnsignedImmediate(args[2])
		if err != nil {
			return 0, err
		}

		instruction = arch.ATypeImmInstruction{
			Opcode:    opcode,
			RegDest:   parsedRegs[0],
			RegA:      parsedRegs[1],
			Immediate: imm,
		}
	}

	return arch.EncodeATypeImmInstruction(instruction), nil
}

func assembleMTypeInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	if len(args) != 3 {
		return 0, ErrInvalidArgumentCount
	}

	parsedRegs, err := parseRegisters(args[0:2])
	if err != nil {
		return 0, err
	}

	imm, err := parseSignedImmediate(args[2])
	if err != nil {
		return 0, err
	}

	instruction := arch.MTypeInstruction{
		Opcode:    opcode,
		RegA:      parsedRegs[0],
		RegB:      parsedRegs[1],
		Immediate: imm,
	}
	return arch.EncodeMTypeInstruction(instruction), nil
}

func assembleETypeInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	if len(args) != 2 {
		return 0, ErrInvalidArgumentCount
	}

	regDest, err := parseRegister(args[0])
	if err != nil {
		return 0, err
	}

	imm, err := parseUnsignedImmediate(args[1])
	if err != nil {
		return 0, err
	}

	instruction := arch.ETypeInstruction{
		Opcode:    opcode,
		RegDest:   regDest,
		Immediate: imm,
	}
	return arch.EncodeETypeInstruction(instruction), nil
}

func assembleBTypeInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	if len(args) != 1 {
		return 0, ErrInvalidArgumentCount
	}

	regA, err := parseRegister(args[0])
	if err != nil {
		return 0, err
	}

	instruction := arch.BTypeInstruction{
		Opcode: opcode,
		RegA:   regA,
	}
	return arch.EncodeBTypeInstruction(instruction), nil
}

func assembleBTypeImmInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	if len(args) != 1 {
		return 0, ErrInvalidArgumentCount
	}

	offset, err := parseSignedImmediate(args[0])
	if err != nil {
		return 0, err
	}

	instruction := arch.BTypeImmInstruction{
		Opcode: opcode,
		Offset: offset,
	}
	return arch.EncodeBTypeImmInstruction(instruction), nil
}

func assembleOTypeInstruction(opcode arch.Opcode, args []string) (arch.Instruction, error) {
	if len(args) != 0 {
		return 0, ErrInvalidArgumentCount
	}

	instruction := arch.OTypeInstruction{
		Opcode: opcode,
	}
	return arch.EncodeOTypeInstruction(instruction), nil
}

func parseRegister(registerName string) (arch.RegisterValue, error) {
	if registerName[0] == 'r' || registerName[0] == 'R' {
		registerName = registerName[1:]
	}

	val, err := strconv.ParseUint(registerName, 10, 8)
	if err != nil {
		return 0, err
	}

	return uint8(val), nil
}

func parseRegisters(registers []string) ([]arch.RegisterValue, error) {
	res := make([]arch.RegisterValue, len(registers))
	for i := range registers {
		parsed, err := parseRegister(registers[i])
		if err != nil {
			return nil, err
		}
		res[i] = parsed
	}
	return res, nil
}

func parseUnsignedImmediate(imm string) (uint16, error) {
	if len(imm) < 2 {
		return 0, fmt.Errorf("invalid immediate %q", imm)
	}

	var val uint16
	kind := imm[0]
	rest := imm[1:]

	if kind == '#' { // decimal
		res, err := strconv.ParseUint(rest, 10, 16)
		if err != nil {
			return 0, err
		}
		val = uint16(res)
	} else if kind == 'x' { // hexadecimal
		res, err := strconv.ParseUint(rest, 16, 16)
		if err != nil {
			return 0, err
		}
		val = uint16(res)
	} else if kind == 'o' { // octal
		res, err := strconv.ParseUint(rest, 8, 16)
		if err != nil {
			return 0, err
		}
		val = uint16(res)
	} else {
		return 0, fmt.Errorf("invalid immediate type specifier %q", kind)
	}

	return val, nil
}

func parseSignedImmediate(imm string) (int16, error) {
	if len(imm) < 2 {
		return 0, fmt.Errorf("invalid immediate %q", imm)
	}

	var val int16
	kind := imm[0]
	rest := imm[1:]

	if kind == '#' { // decimal
		res, err := strconv.ParseInt(rest, 10, 16)
		if err != nil {
			return 0, err
		}
		val = int16(res)
	} else if kind == 'x' { // hexadecimal
		res, err := strconv.ParseInt(rest, 16, 16)
		if err != nil {
			return 0, err
		}
		val = int16(res)
	} else if kind == 'o' { // octal
		res, err := strconv.ParseInt(rest, 8, 16)
		if err != nil {
			return 0, err
		}
		val = int16(res)
	} else {
		return 0, fmt.Errorf("invalid immediate type specifier %q", kind)
	}

	return val, nil
}
