package asm

import (
	"fmt"
	"github.com/dnsge/orange/arch"
	"github.com/dnsge/orange/asm/lexer"
	"github.com/dnsge/orange/asm/parser"
)

func assembleDataDirective(s *parser.Statement, state TraversalState) ([]arch.Instruction, error) {
	directiveToken := s.Body[0]
	if directiveToken.Kind == lexer.FILL_STATEMENT {
		var val int64
		if s.Body[1].Kind == lexer.ADDRESS_OF {
			// fill address of a label
			addr, err := determineAddressOf(s.Body[2], state)
			if err != nil {
				return nil, err
			}
			val = int64(addr)
		} else {
			// fill a 64bit immediate
			imm, err := parseSigned64Immediate(s.Body[1])
			if err != nil {
				return nil, err
			}
			val = imm
		}

		return []arch.Instruction{
			arch.Instruction(val & 0xFFFFFFFF),
			arch.Instruction((val >> 32) & 0xFFFFFFFF),
		}, nil
	} else if directiveToken.Kind == lexer.STRING_STATEMENT {
		str := s.Body[1].Value
		return convertStringToWords(str), nil
	} else {
		return nil, fmt.Errorf("assembleDataDirective: unimplemented for directive %v", directiveToken.Kind)
	}
}

func roundUpToMultiple(num, multiple int) int {
	if multiple == 0 {
		return num
	}

	remainder := num % multiple
	if remainder == 0 {
		return num
	}

	return num + multiple - remainder
}

func calculateStringByteCount(str string) int {
	strBytes := len(str) + 1
	strPaddedBytes := roundUpToMultiple(strBytes, 4)
	return strPaddedBytes
}

func convertStringToWords(str string) []arch.Instruction {
	asBytes := []byte(str)
	asBytes = append(asBytes, 0) // add null terminator
	numWords := roundUpToMultiple(len(asBytes), 4) / 4

	var allWords []arch.Instruction
	for i := 0; i < numWords; i++ {
		start := i * 4
		end := start + 4
		part := []byte{0, 0, 0, 0}
		if end > len(asBytes) {
			end = len(asBytes)
		}
		copy(part, asBytes[start:end])
		packed := arch.ByteOrder.Uint32(part)
		allWords = append(allWords, packed)
	}

	return allWords
}
