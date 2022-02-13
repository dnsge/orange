package lexer

import "fmt"

// Expect returns an ExpectationEntry that captures an expected TokenKind
func Expect(kind TokenKind) ExpectationEntry {
	return &singleExpectationEntry{
		kind: kind,
		keep: true,
	}
}

// ExpectIgnore returns an ExpectationEntry that expects but ignores a TokenKind
func ExpectIgnore(kind TokenKind) ExpectationEntry {
	return &singleExpectationEntry{
		kind: kind,
		keep: false,
	}
}

// ExpectAny returns an ExpectationEntry that captures one of any given TokenKind
func ExpectAny(kinds ...TokenKind) ExpectationEntry {
	return &multipleExpectationEntry{
		kinds: kinds,
		keep:  true,
	}
}

// ExpectAnyIgnore returns an ExpectationEntry that expects but ignores one of any given TokenKind
func ExpectAnyIgnore(kinds ...TokenKind) ExpectationEntry {
	return &multipleExpectationEntry{
		kinds: kinds,
		keep:  false,
	}
}

type singleExpectationEntry struct {
	kind TokenKind
	keep bool
}

func (s *singleExpectationEntry) Matches(kind TokenKind) bool {
	return s.kind == kind
}

func (s *singleExpectationEntry) Keep() bool {
	return s.keep
}

func (s *singleExpectationEntry) Describe() string {
	return DescribeTokenKind(s.kind)
}

type multipleExpectationEntry struct {
	kinds []TokenKind
	keep  bool
}

func (m *multipleExpectationEntry) Matches(kind TokenKind) bool {
	for _, k := range m.kinds {
		if k == kind {
			return true
		}
	}
	return false
}

func (m *multipleExpectationEntry) Keep() bool {
	return m.keep
}

func (m *multipleExpectationEntry) Describe() string {
	if len(m.kinds) == 1 {
		return DescribeTokenKind(m.kinds[0])
	}

	res := "one of "
	for i := range m.kinds {
		res += DescribeTokenKind(m.kinds[i])
		if i != len(m.kinds)-1 {
			res += ", "
		}
	}

	return res
}

// ExpectationEntry describes a possible token expectation that can be matched
// within stream of many Tokens.
type ExpectationEntry interface {
	Matches(kind TokenKind) bool
	Keep() bool
	Describe() string
}

// Expectation is an aggregate of ExpectationEntries that matches multiple
// expected Tokens.
type Expectation struct {
	keepCount int
	entries   []ExpectationEntry
}

func NewExpectation(entries ...ExpectationEntry) *Expectation {
	keepCount := 0
	for _, entry := range entries {
		if entry.Keep() {
			keepCount++
		}
	}

	return &Expectation{
		keepCount: keepCount,
		entries:   entries,
	}
}

// ExtractionCount returns the number of captured (e.g. not ignored) Tokens
// that its pattern describes
func (e *Expectation) ExtractionCount() int {
	return e.keepCount
}

// ExtractExpectedStructure attempts to extract a subset of Tokens from a
// TokenStream, storing the results in dest. dest must have a size large
// enough to store all the expected tokens. offset is the index offset for
// Tokens being stored in dest.
func ExtractExpectedStructure(stream *TokenStream, dest []*Token, exp *Expectation, offset int) error {
	i := 0
	for _, e := range exp.entries {
		actual := stream.Pop()
		if e.Matches(actual.Kind) {
			if e.Keep() {
				dest[offset+i] = actual
				i++
			}
		} else {
			return fmt.Errorf("unexpected token %s at %d:%d (expected %s)", DescribeToken(actual), actual.Row, actual.Column, e.Describe())
		}
	}
	return nil
}
