package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
)

var (
	ErrInvalidOpcode = fmt.Errorf("invalid opcode")
)

func (a *assemblyContext) assembleInstruction(opStatement *parser.Statement) (arch.Instruction, error) {
	opcodeToken := opStatement.Body[0]
	args := opStatement.Body[1:]

	opcode := lexer.GetTokenOpOpcode(opcodeToken.Kind)
	iType := arch.GetInstructionType(opcode)
	switch iType {
	case arch.IType_A:
		return assembleATypeInstruction(opcode, args)
	case arch.IType_AI:
		return assembleATypeImmInstruction(opcode, args)
	case arch.IType_M:
		return assembleMTypeInstruction(opcode, args)
	case arch.IType_E:
		return assembleETypeInstruction(opcode, args)
	case arch.IType_BI:
		return assembleBTypeImmInstruction(opcode, args, a)
	case arch.IType_B:
		return assembleBTypeInstruction(opcode, args)
	case arch.IType_R:
		return assembleRTypeInstruction(opcode, args)
	case arch.IType_O:
		return assembleOTypeInstruction(opcode, args)
	default:
		return 0, ErrInvalidOpcode
	}
}
