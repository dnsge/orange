package asm

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"io"
	"math"
)

type labelMap map[string]uint32

type assemblyContext struct {
	statements  []*parser.Statement
	labels      labelMap
	currAddress uint32
}

func (a *assemblyContext) AddressFor(tok *lexer.Token) (uint32, error) {
	labelAddr, ok := a.labels[tok.Value]
	if !ok {
		return 0, fmt.Errorf("undefined label %q at %d:%d", tok.Value, tok.Row, tok.Column)
	}
	return labelAddr, nil
}

func (a *assemblyContext) OffsetFor(tok *lexer.Token) (uint16, error) {
	labelAddr, ok := a.labels[tok.Value]
	if !ok {
		return 0, fmt.Errorf("undefined label %q at %d:%d", tok.Value, tok.Row, tok.Column)
	}

	computed := int32(labelAddr) - int32(a.currAddress)
	if computed > math.MaxUint16 || computed < 0 {
		return 0, fmt.Errorf("cannot represent label %s with relative with offset %d at %d:%d", tok.Value, computed, tok.Row, tok.Column)
	}
	return uint16(computed), nil
}

func (a *assemblyContext) SignedOffsetFor(tok *lexer.Token) (int16, error) {
	labelAddr, ok := a.labels[tok.Value]
	if !ok {
		return 0, fmt.Errorf("undefined label %q at %d:%d", tok.Value, tok.Row, tok.Column)
	}

	computed := int32(labelAddr) - int32(a.currAddress)
	if computed > math.MaxInt16 || computed < math.MinInt16 {
		return 0, fmt.Errorf("cannot represent label %s with relative with offset %d at %d:%d", tok.Value, computed, tok.Row, tok.Column)
	}
	return int16(computed), nil
}

func AssembleStream(inputFile io.Reader, outputFile io.Writer) error {
	aContext := &assemblyContext{
		labels:      make(labelMap),
		currAddress: 0,
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
			if s.Relocate != nil {
				err = s.Relocate(aContext)
				if err != nil {
					return err
				}
			}
			printStatement(s)

			if assembled, err := aContext.assembleInstruction(s); err != nil {
				return err
			} else if err = binary.Write(outputFile, binary.LittleEndian, assembled); err != nil {
				return err
			}
			aContext.currAddress += 4
		} else if s.Kind == parser.DirectiveStatement && IsDataDirective(s.Body[0].Kind) {
			assembled, err := aContext.assembleDataDirective(s)
			if err != nil {
				return err
			}

			printStatement(s)
			for i := range assembled {
				if err = binary.Write(outputFile, binary.LittleEndian, assembled[i]); err != nil {
					return err
				}
				aContext.currAddress += 4
			}
		}
	}

	return nil
}

// IsDataDirective returns whether the TokenKind represents a directive that
// will appear as data in the final assembled binary
func IsDataDirective(kind lexer.TokenKind) bool {
	return kind == lexer.FILL_STATEMENT || kind == lexer.STRING_STATEMENT
}

func printStatement(statement *parser.Statement) {
	for _, tok := range statement.Body {
		fmt.Printf("%s ", lexer.DescribeToken(tok))
	}
	fmt.Printf("\n")
}
