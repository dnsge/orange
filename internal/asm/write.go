package asm

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"io"
)

type assemblyContext struct {
	statements []*parser.Statement
}

func AssembleStream(inputFile io.Reader, outputFile io.Writer) error {
	aContext := &assemblyContext{}

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

	layout := newLayout()
	err = layout.InitWithStatements(aContext.statements)
	if err != nil {
		return err
	}

	err = layout.Traverse(func(s *parser.Statement, state TraversalState) error {
		if s.Kind == parser.InstructionStatement {
			if s.Relocate != nil {
				err = s.Relocate(state)
				if err != nil {
					return err
				}
			}
			printStatement(s)

			if assembled, err := aContext.assembleInstruction(s, state); err != nil {
				return err
			} else if err = binary.Write(outputFile, binary.LittleEndian, assembled); err != nil {
				return err
			}
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
			}
		}

		return nil
	})

	return err
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
