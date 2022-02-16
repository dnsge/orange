package parser

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"math"
)

var pseudoStatements = []pseudoStatement{
	&opcodePseudoStatement{
		opcode: lexer.CMP,
		convert: func(cmpStatement *Statement) ([]*Statement, error) {
			// CMP r1, r2
			// will become
			// SUB r0, r1, r2

			newBody := []*lexer.Token{
				remapToken(cmpStatement.Body[0], lexer.SUB, "SUB"),
				blankToken(lexer.REGISTER, "r0"),
				cmpStatement.Body[1],
				cmpStatement.Body[2],
			}

			return []*Statement{{
				Body: newBody,
				Kind: InstructionStatement,
			}}, nil
		},
	},
	&opcodePseudoStatement{
		opcode: lexer.CMPI,
		convert: func(cmpiStatement *Statement) ([]*Statement, error) {
			// CMPI r1, #imm
			// will become
			// SUBI r0, r1, #imm

			newBody := []*lexer.Token{
				remapToken(cmpiStatement.Body[0], lexer.SUB, "SUBI"),
				blankToken(lexer.REGISTER, "r0"),
				cmpiStatement.Body[1],
				cmpiStatement.Body[2],
			}

			return []*Statement{{
				Body: newBody,
				Kind: InstructionStatement,
			}}, nil
		},
	},
	&opcodePseudoStatement{
		opcode: lexer.ADR,
		convert: func(adrStatement *Statement) ([]*Statement, error) {
			// ADR r1, $label
			// will become
			// MOVZ r1, #absAddress

			movStatement := &Statement{
				Body: []*lexer.Token{
					remapToken(adrStatement.Body[0], lexer.MOVZ, "MOVZ"),
					adrStatement.Body[1],
					remapToken(adrStatement.Body[1], lexer.BASE_10_IMM, "#0"),
				}, Kind: InstructionStatement,
			}

			movStatement.Relocate = func(relocator Relocator) error {
				address, err := relocator.AddressFor(adrStatement.Body[2])
				if err != nil {
					return err
				}

				if address > math.MaxUint16 {
					return fmt.Errorf("unable to represent absolute address %d in 16 bits", address)
				}
				movStatement.Body[2].Value = fmt.Sprintf("#%d", address)
				return nil
			}

			return []*Statement{movStatement}, nil
		},
	},
}

func translateStatement(opStatement *Statement) ([]*Statement, error) {
	for i := range pseudoStatements {
		if pseudoStatements[i].Matches(opStatement) {
			return pseudoStatements[i].Convert(opStatement)
		}
	}
	return []*Statement{opStatement}, nil
}

func remapToken(tok *lexer.Token, kind lexer.TokenKind, value string) *lexer.Token {
	return &lexer.Token{
		Kind:   kind,
		Value:  value,
		Row:    tok.Row,
		Column: tok.Column,
	}
}

func blankToken(kind lexer.TokenKind, value string) *lexer.Token {
	return &lexer.Token{
		Kind:   kind,
		Value:  value,
		Row:    -1,
		Column: -1,
	}
}
