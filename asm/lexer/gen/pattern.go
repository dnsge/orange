package main

var (
	DefaultPattern = Pattern{nil}
)

type Pattern struct {
	LexerPatterns []string
}

func (p *Pattern) IsDefault() bool {
	return p.LexerPatterns == nil
}

func Only(pattern string) Pattern {
	return Pattern{[]string{pattern}}
}

func OneOf(patterns ...string) Pattern {
	return Pattern{patterns}
}
