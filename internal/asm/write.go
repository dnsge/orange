package asm

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"io"
)

func readFileAndLayout(inputFile io.Reader) (*Layout, error) {
	rawData, err := io.ReadAll(inputFile)
	if err != nil {
		return nil, err
	}

	// tokenize the raw file
	tokens, err := parser.TokenizeAll(rawData)
	if err != nil {
		return nil, err
	}

	// parse tokens into statements
	statements, err := parser.ParseTokens(tokens)
	if err != nil {
		return nil, err
	}

	// initialize the layout
	layout := newLayout()
	err = layout.InitWithStatements(statements)
	if err != nil {
		return nil, err
	}

	return layout, nil
}

func AssembleExecutable(inputFile io.Reader, outputFile io.Writer) error {
	layout, err := readFileAndLayout(inputFile)
	if err != nil {
		return err
	}

	err = layout.Assemble(AssembleStatement)
	if err != nil {
		return err
	}

	err = layout.Traverse(func(section *Section) error {
		for _, a := range section.AssembledStatements {
			err = binary.Write(outputFile, byteOrder, a)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// AssembleStatement turns a parser.Statement into a 32-bit word that will exist in
// the final binary for a program.
//
// Instructions are assembled according to the ISA and data directives, like .fill
// or .string are assembled to include the raw data in the form of 32-bit words.
//
// Strings are thus padded with null bytes at the end to make the string occupy a
// multiple of 32 bits.
func AssembleStatement(s *parser.Statement, state TraversalState) ([]arch.Instruction, error) {
	if s.Relocate != nil {
		err := s.Relocate(state)
		if err != nil {
			return nil, err
		}
	}

	printStatement(s)
	if s.Kind == parser.InstructionStatement {
		assembled, err := assembleInstruction(s, state)
		if err != nil {
			return nil, err
		}
		return []arch.Instruction{assembled}, nil
	} else if s.Kind == parser.DirectiveStatement && IsDataDirective(s.Body[0].Kind) {
		return assembleDataDirective(s)
	} else {
		return []arch.Instruction{}, nil
	}
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
