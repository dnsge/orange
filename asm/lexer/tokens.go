//go:generate go run github.com/dnsge/orange/asm/lexer/gen ./generated_tokens.go
package lexer

// Token describes a lexeme within an input
type Token struct {
	Kind   TokenKind
	Value  string
	Row    int
	Column int
}
