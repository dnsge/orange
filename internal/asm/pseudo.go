package asm

import "github.com/dnsge/orange/internal/asm/lexer"

type pseudoStatement interface {
	Matches(s *statement) bool
	Convert(s *statement) ([]*statement, error)
}

type opcodePseudoStatement struct {
	opcode  lexer.TokenKind
	convert func(s *statement) ([]*statement, error)
}

func (ops *opcodePseudoStatement) Matches(s *statement) bool {
	if len(s.body) == 0 {
		return false
	}
	return s.body[0].Kind == ops.opcode
}

func (ops *opcodePseudoStatement) Convert(s *statement) ([]*statement, error) {
	return ops.convert(s)
}
