package parser

import "github.com/dnsge/orange/internal/asm/lexer"

var pseudoStatements = []pseudoStatement{
	&opcodePseudoStatement{
		opcode:  lexer.CMP,
		convert: translateCMP,
	},
	&opcodePseudoStatement{
		opcode:  lexer.CMPI,
		convert: translateCMPI,
	},
	&opcodePseudoStatement{
		opcode:  lexer.ADR,
		convert: translateADR,
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

func translateCMP(cmpStatement *Statement) ([]*Statement, error) {
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
}

func translateCMPI(cmpiStatement *Statement) ([]*Statement, error) {
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
}

func translateADR(adrStatement *Statement) ([]*Statement, error) {
	newBody := []*lexer.Token{
		remapToken(adrStatement.Body[0], lexer.LDREG, "LDREG"),
		adrStatement.Body[1],
		blankToken(lexer.REGISTER, "r0"),
		blankToken(lexer.BASE_10_IMM, "#0"),
	}

	return []*Statement{{
		Body: newBody,
		Kind: InstructionStatement,
	}}, nil
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
