package asm

import (
	"github.com/dnsge/orange/internal/asm/lexer"
	"io"
)

type labelMap map[string]uint32

type assemblyContext struct {
	tokens   []*lexer.Token
	labels   labelMap
	currLine uint32
}

func AssembleStream(inputFile io.Reader, outputFile io.Writer) error {
	aContext := assemblyContext{
		tokens:   nil,
		labels:   make(labelMap),
		currLine: 0,
	}

	rawData, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	// tokenize the raw file
	if err := aContext.tokenizeAll(rawData); err != nil {
		return err
	}

	// construct statements (instruction + directive) from the generated tokens
	if err := aContext.parseTokens(); err != nil {
		return err
	}

	return nil
}
