package asm

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
}

func translateStatement(opStatement *statement) ([]*statement, error) {
	for i := range pseudoStatements {
		if pseudoStatements[i].Matches(opStatement) {
			return pseudoStatements[i].Convert(opStatement)
		}
	}
	return []*statement{opStatement}, nil
}

func translateCMP(cmpStatement *statement) ([]*statement, error) {
	newBody := []*lexer.Token{
		remapToken(cmpStatement.body[0], lexer.SUB, "SUB"),
		blankToken(lexer.REGISTER, "r0"),
		cmpStatement.body[1],
		cmpStatement.body[2],
	}

	return []*statement{{
		body: newBody,
		kind: instructionStatement,
	}}, nil
}

func translateCMPI(cmpiStatement *statement) ([]*statement, error) {
	newBody := []*lexer.Token{
		remapToken(cmpiStatement.body[0], lexer.SUB, "SUBI"),
		blankToken(lexer.REGISTER, "r0"),
		cmpiStatement.body[1],
		cmpiStatement.body[2],
	}

	return []*statement{{
		body: newBody,
		kind: instructionStatement,
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
