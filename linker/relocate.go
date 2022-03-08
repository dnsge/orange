package linker

import (
	"fmt"
	"github.com/dnsge/orange/arch"
	"github.com/dnsge/orange/asm/lexer"
	"github.com/dnsge/orange/linker/objfile"
	"math"
)

// performRelocation actually completes the task of relocating a specific
// instruction given the target addresses and the relocation entry.
//
// Currently, only the following instructions are relocated:
//  - .fill .addressOf $label
//  - B.EQ $label
//  - MOVZ r1, .addressOf $label / ADR r1, $label
func performRelocation(instructions []arch.Instruction, symbolAddress int, relocateAddress int, relocation *objfile.RelocationTableEntry) error {
	target := &instructions[relocateAddress/4]
	switch relocation.StatementToken {
	case lexer.FILL_STATEMENT:
		// In the case of fill, we must have filled the address with for symbol.
		// Therefore, we can simply replace the value with the absolute address
		// for the target symbol.
		*target = arch.Instruction(symbolAddress)
		return nil
	}

	if lexer.IsTokenOp(relocation.StatementToken) {
		opcode := lexer.GetTokenOpOpcode(relocation.StatementToken)
		kind := arch.GetInstructionType(opcode)
		switch kind {
		case arch.IType_BI:
			// Handle relative branches like B, BL, B.EQ, etc.
			bImmInstruction := arch.DecodeBTypeImmInstruction(*target, opcode)
			if offset, err := computeInstructionOffset(symbolAddress, relocateAddress); err != nil {
				return err
			} else {
				bImmInstruction.Offset = offset
				*target = arch.EncodeBTypeImmInstruction(bImmInstruction)
				return nil
			}
		case arch.IType_E:
			// Handle MOVZ, MOVK
			eInstruction := arch.DecodeETypeInstruction(*target, opcode)
			if address, err := convertAddressToImmediate(symbolAddress); err != nil {
				return err
			} else {
				eInstruction.Immediate = address
				*target = arch.EncodeETypeInstruction(eInstruction)
				return nil
			}
		}
	}

	return fmt.Errorf("unable to perform relocation for token %s", lexer.DescribeTokenKind(relocation.StatementToken))
}

func computeInstructionOffset(target, current int) (int16, error) {
	instructionDiff := (target - current) / 4
	if instructionDiff > math.MaxInt16 || instructionDiff < math.MinInt16 {
		return 0, fmt.Errorf("cannot represent relocated offset %d in int16", instructionDiff)
	}
	return int16(instructionDiff), nil
}

func convertAddressToImmediate(target int) (uint16, error) {
	instructionDiff := target
	if instructionDiff > math.MaxUint16 || instructionDiff < 0 {
		return 0, fmt.Errorf("cannot represent relocated address %d in uint16", instructionDiff)
	}
	return uint16(instructionDiff), nil
}
