package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"strings"
)

var (
	ErrEmptyInstruction = fmt.Errorf("empty instruction")
	ErrInvalidOpcode    = fmt.Errorf("invalid opcode")
)

func ParseAssembly(line string) (arch.Instruction, error) {
	if line == "" {
		return 0, ErrEmptyInstruction
	}

	line = strings.ReplaceAll(line, ",", "")
	tokens := strings.Split(line, " ")

	opcodeText := tokens[0]
	opcode, err := parseOpcodeText(opcodeText)
	if err != nil {
		return 0, err
	}

	var args []string
	if len(tokens) == 1 {
		args = []string{}
	} else {
		args = tokens[1:]
	}

	switch opcode {
	case arch.ADD,
		arch.SUB,
		arch.AND,
		arch.OR,
		arch.XOR,
		arch.LSL,
		arch.LSR,
		arch.CMP:
		return assembleATypeInstruction(opcode, args)
	case arch.ADDI,
		arch.SUBI,
		arch.CMPI:
		return assembleATypeImmInstruction(opcode, args)
	case arch.LDREG,
		arch.LDWORD,
		arch.LDHWRD,
		arch.LDBYTE,
		arch.STREG,
		arch.STWORD,
		arch.STHWRD,
		arch.STBYTE:
		return assembleMTypeInstruction(opcode, args)
	case arch.MOVZ,
		arch.MOVK:
		return assembleETypeInstruction(opcode, args)
	case arch.B,
		arch.B_EQ,
		arch.B_NEQ,
		arch.B_LT,
		arch.B_LE,
		arch.B_GT,
		arch.B_GE:
		return assembleBTypeImmInstruction(opcode, args)
	case arch.BREG:
		return assembleBTypeInstruction(opcode, args)
	case arch.HALT,
		arch.NOOP:
		return assembleOTypeInstruction(opcode, args)
	default:
		return 0, ErrInvalidOpcode
	}
}

func parseOpcodeText(text string) (arch.Opcode, error) {
	switch text {
	case "ADD":
		return arch.ADD, nil
	case "ADDI":
		return arch.ADDI, nil
	case "SUB":
		return arch.SUB, nil
	case "SUBI":
		return arch.SUBI, nil
	case "AND":
		return arch.AND, nil
	case "OR":
		return arch.OR, nil
	case "XOR":
		return arch.XOR, nil
	case "LSL":
		return arch.LSL, nil
	case "LSR":
		return arch.LSR, nil
	case "CMP":
		return arch.CMP, nil
	case "CMPI":
		return arch.CMPI, nil
	case "LDREG":
		return arch.LDREG, nil
	case "LDWORD":
		return arch.LDWORD, nil
	case "LDHWRD":
		return arch.LDHWRD, nil
	case "LDBYTE":
		return arch.LDBYTE, nil
	case "STREG":
		return arch.STREG, nil
	case "STWORD":
		return arch.STWORD, nil
	case "STHWRD":
		return arch.STHWRD, nil
	case "STBYTE":
		return arch.STBYTE, nil
	case "MOVZ":
		return arch.MOVZ, nil
	case "MOVK":
		return arch.MOVK, nil
	case "B":
		return arch.B, nil
	case "BREG":
		return arch.BREG, nil
	case "B.EQ":
		return arch.B_EQ, nil
	case "B.NEQ":
		return arch.B_NEQ, nil
	case "B.LT":
		return arch.B_LT, nil
	case "B.LE":
		return arch.B_LE, nil
	case "B.GT":
		return arch.B_GT, nil
	case "B.GE":
		return arch.B_GE, nil
	case "BL":
		return arch.BL, nil
	case "HALT":
		return arch.HALT, nil
	case "NOOP":
		return arch.NOOP, nil
	default:
		return arch.None, ErrInvalidOpcode
	}
}
