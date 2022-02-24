package asm

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
)

var byteOrder = binary.LittleEndian

func assembleDataDirective(s *parser.Statement) ([]arch.Instruction, error) {
	directiveToken := s.Body[0]
	if directiveToken.Kind == lexer.FILL_STATEMENT {
		// fill a 64bit immediate
		val, err := parseSigned64Immediate(s.Body[1])
		if err != nil {
			return nil, err
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
		packed := byteOrder.Uint32(part)
		allWords = append(allWords, packed)
	}

	return allWords
}
