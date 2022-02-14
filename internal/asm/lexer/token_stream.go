package lexer

// NewTokenStream returns a TokenStream encapsulating tokens.
//
// tokens should not be modified while TokenStream is in use.
func NewTokenStream(tokens []*Token) *TokenStream {
	return &TokenStream{
		pos:    0,
		tokens: tokens,
	}
}

// TokenStream describes a stream of lexer Tokens
type TokenStream struct {
	pos    int
	tokens []*Token
}

// HasNext returns whether the stream has more tokens
func (ts *TokenStream) HasNext() bool {
	return ts.pos < len(ts.tokens)
}

// Remaining returns how many tokens remain in the stream
func (ts *TokenStream) Remaining() int {
	return len(ts.tokens) - ts.pos
}

// Peek returns the next Token without removing it
func (ts *TokenStream) Peek() *Token {
	return ts.tokens[ts.pos]
}

// Pop removes and returns the next Token, advancing the stream
func (ts *TokenStream) Pop() *Token {
	if ts.pos >= len(ts.tokens) {
		return nil
	}
	oldPos := ts.pos
	ts.pos++
	return ts.tokens[oldPos]
}
