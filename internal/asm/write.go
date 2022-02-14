package asm

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"io"
)

type labelMap map[string]uint32

type assemblyContext struct {
	statements []*parser.Statement
	labels     labelMap
	currLine   uint32
}

func AssembleStream(inputFile io.Reader, outputFile io.Writer) error {
	aContext := assemblyContext{
		labels:   make(labelMap),
		currLine: 0,
	}

	rawData, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	// tokenize the raw file
	tokens, err := parser.TokenizeAll(rawData)
	if err != nil {
		return err
	}

	// parse tokens into statements
	statements, err := parser.ParseTokens(tokens)
	if err != nil {
		return err
	}

	aContext.statements = statements
	if err := aContext.processLabelDeclarations(); err != nil {
		return err
	}

	fmt.Printf("%#+v\n", aContext.labels)

	for _, s := range aContext.statements {
		if s.Kind == parser.InstructionStatement {
			if assembled, err := aContext.assembleInstruction(s); err != nil {
				return err
			} else if err = binary.Write(outputFile, binary.LittleEndian, assembled); err != nil {
				return err
			}
			aContext.currLine++
		} else if s.Kind == parser.DirectiveStatement && lexer.IsDataDirective(s.Body[0].Kind) {
			assembled, err := aContext.assembleDataDirective(s)
			if err != nil {
				return err
			}

			for i := range assembled {
				if err = binary.Write(outputFile, binary.LittleEndian, assembled[i]); err != nil {
					return err
				}
				aContext.currLine++
			}
		}
	}

	return nil
}
