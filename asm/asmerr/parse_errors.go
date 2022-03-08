package asmerr

import (
	"fmt"
	"github.com/dnsge/orange/asm/lexer"
)

func describeLocatedToken(token *lexer.Token) string {
	return fmt.Sprintf("%q at %d:%d", lexer.DescribeToken(token), token.Row, token.Column)
}

type InvalidImmediateError struct {
	Token *lexer.Token
}

func (i *InvalidImmediateError) Error() string {
	return fmt.Sprintf("invalid immediate %s", describeLocatedToken(i.Token))
}

type LabelNotFoundError struct {
	Label *lexer.Token
}

func (l *LabelNotFoundError) Error() string {
	return fmt.Sprintf("undefined label %s", describeLocatedToken(l.Label))
}

type DuplicateLabelError struct {
	Label *lexer.Token
	Other *lexer.Token
}

func (d *DuplicateLabelError) Error() string {
	return fmt.Sprintf("duplicate label %s (other: %s)", describeLocatedToken(d.Label), describeLocatedToken(d.Other))
}
