package lexer

func NewTokenStream(tokens []*Token) *TokenStream {
	return &TokenStream{
		pos:    0,
		tokens: tokens,
	}
}

type TokenStream struct {
	pos    int
	tokens []*Token
}

func (ts *TokenStream) HasNext() bool {
	return ts.pos < len(ts.tokens)
}

func (ts *TokenStream) Remaining() int {
	return len(ts.tokens) - ts.pos
}

func (ts *TokenStream) Peek() *Token {
	return ts.tokens[ts.pos]
}

func (ts *TokenStream) Pop() *Token {
	oldPos := ts.pos
	ts.pos++
	return ts.tokens[oldPos]
}
