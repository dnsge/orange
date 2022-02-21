package asm

import (
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/asmerr"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"strconv"
)

var (
	registerAliases = map[string]arch.RegisterValue{
		"rsp": arch.StackRegister,
	}
)

func assembleATypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 3 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 3,
			Got:      len(args),
		}
	}

	parsedRegs, err := parseRegisters(args)
	if err != nil {
		return 0, err
	}

	instruction := arch.ATypeInstruction{
		Opcode:  opcode,
		RegDest: parsedRegs[0],
		RegA:    parsedRegs[1],
		RegB:    parsedRegs[2],
	}

	return arch.EncodeATypeInstruction(instruction), nil
}

func assembleATypeImmInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 3 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 3,
			Got:      len(args),
		}
	}

	parsedRegs, err := parseRegisters(args[0:2])
	if err != nil {
		return 0, err
	}

	imm, err := parseUnsignedImmediate(args[2])
	if err != nil {
		return 0, err
	}

	instruction := arch.ATypeImmInstruction{
		Opcode:    opcode,
		RegDest:   parsedRegs[0],
		RegA:      parsedRegs[1],
		Immediate: imm,
	}

	return arch.EncodeATypeImmInstruction(instruction), nil
}

func assembleMTypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	var imm int16
	if len(args) == 2 {
		// Handle special no offset case
		// For example, LDREG r3, [r2]
		// aka de-referencing a pointer
		imm = 0
	} else if len(args) == 3 {
		pImm, err := parseSignedImmediate(args[2])
		if err != nil {
			return 0, err
		}
		imm = pImm
	} else {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 3,
			Got:      len(args),
		}
	}

	parsedRegs, err := parseRegisters(args[0:2])
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

func assembleETypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 2 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 2,
			Got:      len(args),
		}
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

func assembleBTypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 1 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 1,
			Got:      len(args),
		}
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

func assembleBTypeImmInstruction(opcode arch.Opcode, args []*lexer.Token, relocator parser.Relocator) (arch.Instruction, error) {
	if len(args) != 1 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 1,
			Got:      len(args),
		}
	}

	offset, err := parseOffsetOrLabel(args[0], relocator)
	if err != nil {
		return 0, err
	}

	instruction := arch.BTypeImmInstruction{
		Opcode: opcode,
		Offset: offset,
	}
	return arch.EncodeBTypeImmInstruction(instruction), nil
}

func assembleOTypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 0 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 0,
			Got:      len(args),
		}
	}

	instruction := arch.OTypeInstruction{
		Opcode: opcode,
	}
	return arch.EncodeOTypeInstruction(instruction), nil
}

func assembleRTypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 1 {
		return 0, &asmerr.InvalidArgumentCountError{
			Opcode:   opcode,
			Expected: 1,
			Got:      len(args),
		}
	}

	regA, err := parseRegister(args[0])
	if err != nil {
		return 0, err
	}

	instruction := arch.RTypeInstruction{
		Opcode: opcode,
		RegA:   regA,
	}
	return arch.EncodeRTypeInstruction(instruction), nil
}

func parseRegister(registerName *lexer.Token) (arch.RegisterValue, error) {
	if alias, ok := registerAliases[registerName.Value]; ok {
		return alias, nil
	}

	regNumberStr := registerName.Value[1:]
	val, err := strconv.ParseUint(regNumberStr, 10, 8)
	if err != nil {
		return 0, err
	}

	return uint8(val), nil
}

func parseRegisters(registers []*lexer.Token) ([]arch.RegisterValue, error) {
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

func parseUnsignedImmediate(immTok *lexer.Token) (uint16, error) {
	imm := immTok.Value[1:]
	var base int

	switch immTok.Kind {
	case lexer.BASE_10_IMM:
		base = 10
	case lexer.BASE_16_IMM:
		base = 16
	case lexer.BASE_8_IMM:
		base = 8
	default:
		return 0, &asmerr.InvalidImmediateError{Token: immTok}
	}

	res, err := strconv.ParseUint(imm, base, 16)
	if err != nil {
		return 0, err
	}
	return uint16(res), nil
}

func parseSignedImmediate(immTok *lexer.Token) (int16, error) {
	imm := immTok.Value[1:]
	var base int

	switch immTok.Kind {
	case lexer.BASE_10_IMM:
		base = 10
	case lexer.BASE_16_IMM:
		base = 16
	case lexer.BASE_8_IMM:
		base = 8
	default:
		return 0, &asmerr.InvalidImmediateError{Token: immTok}
	}

	res, err := strconv.ParseInt(imm, base, 16)
	if err != nil {
		return 0, err
	}
	return int16(res), nil
}

func parseSigned64Immediate(immTok *lexer.Token) (int64, error) {
	imm := immTok.Value[1:]
	var base int

	switch immTok.Kind {
	case lexer.BASE_10_IMM:
		base = 10
	case lexer.BASE_16_IMM:
		base = 16
	case lexer.BASE_8_IMM:
		base = 8
	default:
		return 0, &asmerr.InvalidImmediateError{Token: immTok}
	}

	res, err := strconv.ParseInt(imm, base, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func parseOffsetOrLabel(tok *lexer.Token, relocator parser.Relocator) (int16, error) {
	if tok.Kind == lexer.LABEL {
		instructionAddressOffset, err := relocator.SignedOffsetFor(tok)
		if err != nil {
			return 0, err
		}

		instructionOffset := instructionAddressOffset / 4
		return instructionOffset, nil
	} else {
		return parseSignedImmediate(tok)
	}
}
