package parser

import (
	"errors"
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/timtadh/lexmachine/machines"
	"io"
)

type StatementKind uint8

const (
	InstructionStatement StatementKind = iota
	DirectiveStatement
)

var (
	ErrExpectedLineEnd = fmt.Errorf("expected line end")
)

// TokenizeAll converts the given data into a slice of Tokens
func TokenizeAll(data []byte) ([]*lexer.Token, error) {
	var allTokens []*lexer.Token
	t, err := lexer.New(data)
	if err != nil {
		return nil, err
	}

	for {
		tok, err := t.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			var unconsumed *machines.UnconsumedInput
			if errors.As(err, &unconsumed) {
				return nil, fmt.Errorf("invalid token at %d:%d", unconsumed.StartLine, unconsumed.StartColumn)

			} else {
				return nil, err
			}
		}

		allTokens = append(allTokens, tok)
	}

	return allTokens, nil
}

// Statement is the most basic element of an Orange assembly program.
//
// It can describe:
//  1. An instruction (e.g. ADD r3, r2, r1)
//  2. A directive (e.g. a label declaration)
type Statement struct {
	Body     []*lexer.Token
	Kind     StatementKind
	Relocate func(relocator Relocator) error
}

type Relocator interface {
	AddressFor(label *lexer.Token) (uint32, bool)
	OffsetFor(label *lexer.Token) (uint16, error)
	SignedOffsetFor(label *lexer.Token) (int16, error)
}

// ParseTokens aggregates tokens from the input stream into statements.
//
// We make the following classification of all complete statements
// within Orange assembly:
//
// All (meaningful) tokens belong to one of two types:
//  1. Instructions
//  2. Directives
//
// Thus, we will iterate over the input token stream and aggregate
// each token to a specific instruction or directive.
//
// For example, given the input token stream state of
//   [MOVZ, REGISTER, COMMA, BASE_10_IMM, COMMENT, LINE_END, ...]
// We will extract MOVZ, inspect the token, determine that the instruction
// should contain a register and an immediate, and remove them from the stream.
//
// The LINE_END token is used to enforce assembly programs to have each instruction
// on a separate line. This is not a limitation of the parsing, but a deliberate
// design choice. Therefore, we must reach a LINE_END token before beginning parsing
// of another Statement.
func ParseTokens(tokens []*lexer.Token) ([]*Statement, error) {
	stream := lexer.NewTokenStream(tokens)

	var statements []*Statement
	waitingForLineEnd := false
	for stream.HasNext() {
		currentToken := stream.Pop()
		var producedStatement *Statement
		if currentToken.Kind == lexer.COMMENT {
			// consume and ignore token
			continue
		} else if currentToken.Kind == lexer.LINE_END {
			waitingForLineEnd = false
			// consume and ignore token
			continue
		} else if lexer.IsTokenOp(currentToken.Kind) {
			// make sure Statement does not trail other Statement
			if waitingForLineEnd {
				return nil, ErrExpectedLineEnd
			}

			if opStatement, err := parseOpTokens(currentToken, stream); err != nil {
				return nil, err
			} else {
				producedStatement = opStatement
				waitingForLineEnd = true
			}
		} else if lexer.IsTokenDirective(currentToken.Kind) {
			// make sure Statement does not trail other Statement
			if waitingForLineEnd {
				return nil, ErrExpectedLineEnd
			}

			if dStatement, err := parseDirectiveTokens(currentToken, stream); err != nil {
				return nil, err
			} else {
				if dStatement.Body[len(dStatement.Body)-1].Kind == lexer.LINE_END {
					// The directive processed the new line, don't update waitingForLineEnd
					dStatement.Body = dStatement.Body[:len(dStatement.Body)-1] // remove last token
				}

				producedStatement = dStatement
			}
		} else {
			return nil, fmt.Errorf("unexpected token %s at %d:%d (expected statement)", lexer.DescribeToken(currentToken), currentToken.Row, currentToken.Column)
		}

		translated, err := translateStatement(producedStatement)
		if err != nil {
			return nil, err
		}

		for i := range translated {
			statements = append(statements, translated[i])
		}
	}

	return statements, nil
}

// parseOpTokens attempts to parse an instruction statement from the TokenStream
func parseOpTokens(op *lexer.Token, stream *lexer.TokenStream) (*Statement, error) {
	exp, err := getOpcodeStatementExpectation(op.Kind)
	if err != nil {
		return nil, err
	}

	statementBody, err := handlePrefixedExtraction(stream, op, exp)
	if err != nil {
		return nil, err
	}

	return &Statement{
		Body: statementBody,
		Kind: InstructionStatement,
	}, nil
}

// parseDirectiveTokens attempts to parse a directive statement from the TokenStream
func parseDirectiveTokens(directive *lexer.Token, stream *lexer.TokenStream) (*Statement, error) {
	exp, ok := directiveKindExpectationMap[directive.Kind]
	if !ok {
		return nil, fmt.Errorf("unexpected directive Kind %v", lexer.DescribeTokenKind(directive.Kind))
	}

	statementBody, err := handlePrefixedExtraction(stream, directive, exp)
	if err != nil {
		return nil, err
	}

	return &Statement{
		Body: statementBody,
		Kind: DirectiveStatement,
	}, nil
}

// handlePrefixedExtraction returns a slice of Tokens after extracting the
// expected tokens from the TokenStream
func handlePrefixedExtraction(stream *lexer.TokenStream, prefix *lexer.Token, extractable lexer.Extractable) ([]*lexer.Token, error) {
	// reserve space for prefix and expectations
	body := make([]*lexer.Token, 1, extractable.ExtractionCount()+1)
	body[0] = prefix
	if err := extractable.Extract(stream, &body); err != nil {
		return nil, err
	} else {
		return body, nil
	}
}
