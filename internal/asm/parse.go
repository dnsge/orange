package asm

import (
	"bufio"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/lexer"
	"io"
	"strings"
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
			return err
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
				statements = append(statements, opStatement)
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

				statements = append(statements, dStatement)
			}
		} else {
			return fmt.Errorf("unexpected token: %s", lexer.DescribeToken(currentToken))
		}
	}

	return nil
}

func parseOpTokens(op *lexer.Token, stream *lexer.TokenStream) (*statement, error) {
	iType := lexer.GetTokenOpInstructionType(op.Kind)
	exp, ok := opKindExpectationMap[iType]
	if !ok {
		return nil, fmt.Errorf("unexpected instruction kind %v", iType)
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

func (a *assemblyContext) processLabels(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line[0] == '.' { // label
			line = strings.TrimSuffix(line, ":")
			labelName := line[1:]
			if _, ok := a.labels[labelName]; ok {
				return fmt.Errorf("duplicate label definition for %q", labelName)
			}
			a.labels[labelName] = a.currLine
		} else {
			a.currLine++
		}
	}

	return nil
}

func (a *assemblyContext) parseAssembly(line string) (arch.Instruction, error) {
	if line == "" {
		return 0, ErrEmptyInstruction
	}

	line = strings.ReplaceAll(line, ",", "")
	tokens := strings.Split(line, " ")

	opcodeText := tokens[0]
	opcode, err := parseOpcodeText(opcodeText)
	if err != nil {
		return 0, err
	}

	var args []string
	if len(tokens) == 1 {
		args = []string{}
	} else {
		args = tokens[1:]
	}

	iType := arch.GetInstructionType(opcode)
	switch iType {
	case arch.IType_A:
		return assembleATypeInstruction(opcode, args)
	case arch.IType_AI:
		return assembleATypeImmInstruction(opcode, args)
	case arch.IType_M:
		return assembleMTypeInstruction(opcode, args)
	case arch.IType_E:
		return assembleETypeInstruction(opcode, args)
	case arch.IType_BI:
		return assembleBTypeImmInstruction(opcode, args, a)
	case arch.IType_B:
		return assembleBTypeInstruction(opcode, args)
	case arch.IType_O:
		return assembleOTypeInstruction(opcode, args)
	default:
		return 0, ErrInvalidOpcode
	}
}

func parseOpcodeText(text string) (arch.Opcode, error) {
	switch text {
	case "ADD":
		return arch.ADD, nil
	case "ADDI":
		return arch.ADDI, nil
	case "SUB":
		return arch.SUB, nil
	case "SUBI":
		return arch.SUBI, nil
	case "AND":
		return arch.AND, nil
	case "OR":
		return arch.OR, nil
	case "XOR":
		return arch.XOR, nil
	case "LSL":
		return arch.LSL, nil
	case "LSR":
		return arch.LSR, nil
	case "CMP":
		return arch.CMP, nil
	case "CMPI":
		return arch.CMPI, nil
	case "LDREG":
		return arch.LDREG, nil
	case "LDWORD":
		return arch.LDWORD, nil
	case "LDHWRD":
		return arch.LDHWRD, nil
	case "LDBYTE":
		return arch.LDBYTE, nil
	case "STREG":
		return arch.STREG, nil
	case "STWORD":
		return arch.STWORD, nil
	case "STHWRD":
		return arch.STHWRD, nil
	case "STBYTE":
		return arch.STBYTE, nil
	case "MOVZ":
		return arch.MOVZ, nil
	case "MOVK":
		return arch.MOVK, nil
	case "B":
		return arch.B, nil
	case "BREG":
		return arch.BREG, nil
	case "B.EQ":
		return arch.B_EQ, nil
	case "B.NEQ":
		return arch.B_NEQ, nil
	case "B.LT":
		return arch.B_LT, nil
	case "B.LE":
		return arch.B_LE, nil
	case "B.GT":
		return arch.B_GT, nil
	case "B.GE":
		return arch.B_GE, nil
	case "BL":
		return arch.BL, nil
	case "HALT":
		return arch.HALT, nil
	case "NOOP":
		return arch.NOOP, nil
	default:
		return arch.None, ErrInvalidOpcode
	}
}
