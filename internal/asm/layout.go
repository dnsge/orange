package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"math"
)

var (
	ErrLabelNotFound = fmt.Errorf("label not found")
)

type Layout struct {
	Sections []*Section
	Labels   map[string]*parser.Statement
}

func newLayout() *Layout {
	return &Layout{
		Sections: []*Section{},
		Labels:   make(map[string]*parser.Statement),
	}
}

type Section struct {
	Name string
	Size int

	Statements     []*parser.Statement
	StatementSizes []int
}

// SectionByName returns the existing section with the given name or returns
// a newly created section with the name at the end of the current binary.
func (l *Layout) SectionByName(name string) *Section {
	for i := range l.Sections {
		if l.Sections[i].Name == name {
			return l.Sections[i]
		}
	}

	// No match, add new section
	newSection := &Section{
		Name: name,
		Size: 0,
	}

	l.Sections = append(l.Sections, newSection)
	return newSection
}

// LocateStatement returns the absolute address of the statement in the final binary
func (l *Layout) LocateStatement(statement *parser.Statement) (int, error) {
	address := 0
	for _, sec := range l.Sections {
		for i := range sec.Statements {
			s := sec.Statements[i]
			if s == statement {
				return address, nil
			}

			sSize := sec.StatementSizes[i]
			address += sSize
		}
	}
	return 0, fmt.Errorf("LocateStatement: not found in any section")
}

// LocateLabel returns the absolute address of the label in the final binary
func (l *Layout) LocateLabel(label string) (uint32, error) {
	labelStatement, ok := l.Labels[label]
	if !ok {
		return 0, fmt.Errorf("label %q: %w", label, ErrLabelNotFound)
	}

	located, err := l.LocateStatement(labelStatement)
	if err != nil {
		return 0, fmt.Errorf("label %q: %w", label, ErrLabelNotFound)
	}

	return uint32(located), nil
}

// InitWithStatements initializes the layout with a list of statements
//
// We assume we are starting with the .text section, regardless of whether
// this function was called earlier and terminated with a different section.
func (l *Layout) InitWithStatements(statements []*parser.Statement) error {
	currentSection := l.SectionByName("text") // initialize text as first section

	for _, s := range statements {
		if s.Kind == parser.DirectiveStatement {
			directiveToken := s.Body[0]
			if directiveToken.Kind == lexer.SECTION {
				// Section begin
				sectionName := s.Body[1].Value
				currentSection = l.SectionByName(sectionName)
			} else if directiveToken.Kind == lexer.LABEL_DECLARATION {
				labelName := directiveToken.Value
				if _, ok := l.Labels[labelName]; ok {
					// duplicate label found
					return fmt.Errorf("duplicate label %q at %d:%d", labelName, directiveToken.Row, directiveToken.Column)
				}

				l.Labels[labelName] = s
			}
		}

		statementSize := CalculateStatementSize(s)
		currentSection.Size += statementSize
		currentSection.Statements = append(currentSection.Statements, s)
		currentSection.StatementSizes = append(currentSection.StatementSizes, statementSize)
	}

	return nil
}

type TraversalState interface {
	parser.Relocator
	Section() *Section
	Address() int
}

// Traverse iterates over each statement throughout the binary, calling the
// traversal function along the way with each statement and a bound TraversalState.
func (l *Layout) Traverse(traversalFunc func(*parser.Statement, TraversalState) error) error {
	address := 0
	for _, sec := range l.Sections {
		for i := range sec.Statements {
			s := sec.Statements[i]
			bound := &boundTraversalState{
				boundLayout:    l,
				section:        sec,
				currentAddress: address,
			}

			// apply traversal function
			err := traversalFunc(s, bound)
			if err != nil {
				return err
			}

			// increment address counter
			sSize := sec.StatementSizes[i]
			address += sSize
		}
	}

	return nil
}

type boundTraversalState struct {
	boundLayout    *Layout
	section        *Section
	currentAddress int
}

func (b *boundTraversalState) Section() *Section {
	return b.section
}

func (b *boundTraversalState) Address() int {
	return b.currentAddress
}

func (b *boundTraversalState) AddressFor(label *lexer.Token) (uint32, error) {
	return b.boundLayout.LocateLabel(label.Value)
}

func (b *boundTraversalState) OffsetFor(label *lexer.Token) (uint16, error) {
	labelAddr, err := b.AddressFor(label)
	if err != nil {
		return 0, fmt.Errorf("undefined label %q at %d:%d", label.Value, label.Row, label.Column)
	}

	computed := int32(labelAddr) - int32(b.currentAddress)
	if computed > math.MaxUint16 || computed < 0 {
		return 0, fmt.Errorf("cannot represent label %s with relative with offset %d at %d:%d", label.Value, computed, label.Row, label.Column)
	}

	return uint16(computed), nil
}

func (b *boundTraversalState) SignedOffsetFor(label *lexer.Token) (int16, error) {
	labelAddr, err := b.AddressFor(label)
	if err != nil {
		return 0, fmt.Errorf("undefined label %q at %d:%d", label.Value, label.Row, label.Column)
	}

	computed := int32(labelAddr) - int32(b.currentAddress)
	if computed > math.MaxInt16 || computed < math.MinInt16 {
		return 0, fmt.Errorf("cannot represent label %s with relative with offset %d at %d:%d", label.Value, computed, label.Row, label.Column)
	}

	return int16(computed), nil
}

// CalculateStatementSize returns the number of bytes that a statement occupies
// within the final binary. For example, most directives take up zero bytes while
// standard instructions take up 4 bytes.
func CalculateStatementSize(statement *parser.Statement) int {
	if statement.Kind == parser.InstructionStatement {
		return 4 // one word per instruction
	} else if statement.Kind == parser.DirectiveStatement {
		directiveToken := statement.Body[0]
		switch directiveToken.Kind {
		case lexer.FILL_STATEMENT:
			return 8 // 8 bytes for register
		case lexer.STRING_STATEMENT:
			str := statement.Body[1].Value
			return calculateStringByteCount(str)
		case lexer.LABEL_DECLARATION, lexer.SECTION:
			return 0
		default:
			return 0
		}
	}

	return 0
}
