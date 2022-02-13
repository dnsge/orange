package parser

import "github.com/dnsge/orange/internal/asm/lexer"

type pseudoStatement interface {
	Matches(s *Statement) bool
	Convert(s *Statement) ([]*Statement, error)
}

type opcodePseudoStatement struct {
	opcode  lexer.TokenKind
	convert func(s *Statement) ([]*Statement, error)
}

func (ops *opcodePseudoStatement) Matches(s *Statement) bool {
	if len(s.Body) == 0 {
		return false
	}
	return s.Body[0].Kind == ops.opcode
}

func (ops *opcodePseudoStatement) Convert(s *Statement) ([]*Statement, error) {
	return ops.convert(s)
}
