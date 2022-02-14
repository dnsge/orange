package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/lexer"
	"math"
	"strconv"
)

var (
	ErrInvalidArgumentCount = fmt.Errorf("invalid argument count")
)

func assembleATypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
	if len(args) != 3 {
		return 0, ErrInvalidArgumentCount
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

	instruction := arch.ATypeImmInstruction{
		Opcode:    opcode,
		RegDest:   parsedRegs[0],
		RegA:      parsedRegs[1],
		Immediate: imm,
	}

	return arch.EncodeATypeImmInstruction(instruction), nil
}

func assembleMTypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
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

func assembleETypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
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

func assembleBTypeInstruction(opcode arch.Opcode, args []*lexer.Token) (arch.Instruction, error) {
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

func assembleBTypeImmInstruction(opcode arch.Opcode, args []*lexer.Token, ctx *assemblyContext) (arch.Instruction, error) {
	if len(args) != 1 {
		return 0, ErrInvalidArgumentCount
	}

	offset, err := parseOffsetOrLabel(args[0], ctx)
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
		return 0, ErrInvalidArgumentCount
	}

	instruction := arch.OTypeInstruction{
		Opcode: opcode,
	}
	return arch.EncodeOTypeInstruction(instruction), nil
}

func parseRegister(registerName *lexer.Token) (arch.RegisterValue, error) {
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
		return 0, fmt.Errorf("invalid immediate type %s", lexer.DescribeTokenKind(immTok.Kind))
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
		return 0, fmt.Errorf("invalid immediate type %s", lexer.DescribeTokenKind(immTok.Kind))
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
		return 0, fmt.Errorf("invalid immediate type %s", lexer.DescribeTokenKind(immTok.Kind))
	}

	res, err := strconv.ParseInt(imm, base, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func parseOffsetOrLabel(tok *lexer.Token, ctx *assemblyContext) (int16, error) {
	if tok.Kind == lexer.LABEL {
		labelName := tok.Value[1:]
		labelTarget, ok := ctx.labels[labelName]
		if !ok {
			return 0, fmt.Errorf("undefined label %q at %d:%d", labelName, tok.Row, tok.Column)
		}

		instructionOffset := int32(labelTarget - ctx.currLine)
		if instructionOffset > math.MaxInt16 || instructionOffset < math.MinInt16 {
			return 0, fmt.Errorf("cannot branch to relative with offset %d (computed at %d:%d)", instructionOffset, tok.Row, tok.Column)
		}

		return int16(instructionOffset), nil
	} else {
		return parseSignedImmediate(tok)
	}
}
