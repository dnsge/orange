package asm

import (
	"errors"
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/timtadh/lexmachine/machines"
	"io"
)

type statementKind uint8

const (
	instructionStatement statementKind = iota
	directiveStatement
)

var (
	ErrEmptyInstruction = fmt.Errorf("empty instruction")
	ErrInvalidOpcode    = fmt.Errorf("invalid opcode")
	ErrExpectedLineEnd  = fmt.Errorf("expected line end")
)

func (a *assemblyContext) tokenizeAll(data []byte) error {
	var allTokens []*lexer.Token
	t, err := lexer.New(data)
	if err != nil {
		return err
	}

	for {
		tok, err := t.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			var unconsumed *machines.UnconsumedInput
			if errors.As(err, &unconsumed) {
				return fmt.Errorf("invalid token at %d:%d", unconsumed.StartLine, unconsumed.StartColumn)

			} else {
				return err
			}
		}

		allTokens = append(allTokens, tok)
	}

	a.tokens = allTokens
	return nil
}

type statement struct {
	body []*lexer.Token
	kind statementKind
}

// parseTokens aggregates tokens from the input stream into statements.
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
// of another statement.
func (a *assemblyContext) parseTokens() error {
	stream := lexer.NewTokenStream(a.tokens)

	var statements []*statement
	waitingForLineEnd := false
	for stream.HasNext() {
		currentToken := stream.Pop()
		var producedStatement *statement
		if currentToken.Kind == lexer.COMMENT {
			// consume and ignore token
			continue
		} else if currentToken.Kind == lexer.LINE_END {
			waitingForLineEnd = false
			// consume and ignore token
			continue
		} else if lexer.IsOp(currentToken.Kind) {
			// make sure statement does not trail other statement
			if waitingForLineEnd {
				return ErrExpectedLineEnd
			}

			if opStatement, err := parseOpTokens(currentToken, stream); err != nil {
				return err
			} else {
				producedStatement = opStatement
				waitingForLineEnd = true
			}
		} else if lexer.IsDirective(currentToken.Kind) {
			// make sure statement does not trail other statement
			if waitingForLineEnd {
				return ErrExpectedLineEnd
			}

			if dStatement, err := parseDirectiveTokens(currentToken, stream); err != nil {
				return err
			} else {
				if dStatement.body[len(dStatement.body)-1].Kind == lexer.LINE_END {
					// The directive processed the new line, don't update waitingForLineEnd
					dStatement.body = dStatement.body[:len(dStatement.body)-1] // remove last token
				} else {
					waitingForLineEnd = true
				}

				producedStatement = dStatement
			}
		} else {
			return fmt.Errorf("unexpected token: %s", lexer.DescribeToken(currentToken))
		}

		translated, err := translateStatement(producedStatement)
		if err != nil {
			return err
		}

		for i := range translated {
			statements = append(statements, translated[i])
			fmt.Printf("%s\n", translated[i].body[0].Value)
		}
	}

	a.statements = statements
	return nil
}

func parseOpTokens(op *lexer.Token, stream *lexer.TokenStream) (*statement, error) {
	exp, err := getOpcodeStatementExpectation(op.Kind)
	if err != nil {
		return nil, err
	}

	statementBody, err := handlePrefixedExtraction(stream, op, exp)
	if err != nil {
		return nil, err
	}

	return &statement{
		body: statementBody,
		kind: instructionStatement,
	}, nil
}

func parseDirectiveTokens(directive *lexer.Token, stream *lexer.TokenStream) (*statement, error) {
	exp, ok := directiveKindExpectationMap[directive.Kind]
	if !ok {
		return nil, fmt.Errorf("unexpected directive kind %v", lexer.DescribeTokenKind(directive.Kind))
	}

	statementBody, err := handlePrefixedExtraction(stream, directive, exp)
	if err != nil {
		return nil, err
	}

	return &statement{
		body: statementBody,
		kind: directiveStatement,
	}, nil
}

func handlePrefixedExtraction(stream *lexer.TokenStream, prefix *lexer.Token, expectation *lexer.Expectation) ([]*lexer.Token, error) {
	body := make([]*lexer.Token, expectation.ExtractionCount()+1)
	body[0] = prefix
	if err := lexer.ExtractExpectedStructure(stream, body, expectation, 1); err != nil {
		return nil, err
	} else {
		return body, nil
	}
}
