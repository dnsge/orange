package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/asmerr"
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

	Statements          []*parser.Statement
	StatementSizes      []int
	AssembledStatements []arch.Instruction
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
func (l *Layout) LocateStatement(statement *parser.Statement) (int, bool) {
	address := 0
	for _, sec := range l.Sections {
		for i := range sec.Statements {
			s := sec.Statements[i]
			if s == statement {
				return address, true
			}

			sSize := sec.StatementSizes[i]
			address += sSize
		}
	}
	return 0, false
}

// LocateLabel returns the absolute address of the label in the final binary
func (l *Layout) LocateLabel(label string) (uint32, bool) {
	labelStatement, ok := l.Labels[label]
	if !ok {
		return 0, false
	}

	located, found := l.LocateStatement(labelStatement)
	if !found {
		return 0, false
	}

	return uint32(located), true
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
				if other, ok := l.Labels[labelName]; ok {
					// duplicate label found
					return &asmerr.DuplicateLabelError{
						Label: directiveToken,
						Other: other.Body[0],
					}
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
func (l *Layout) Traverse(traversalFunc func(section *Section) error) error {
	for _, sec := range l.Sections {
		err := traversalFunc(sec)
		if err != nil {
			return err
		}
	}

	return nil
}

// Assemble iterates over each statement throughout the binary, calling the
// assembler function along the way with each statement and a bound TraversalState.
func (l *Layout) Assemble(assembleFunc func(*parser.Statement, TraversalState) ([]arch.Instruction, error)) error {
	address := 0
	for _, sec := range l.Sections {
		j := 0
		for i := range sec.Statements {
			s := sec.Statements[i]
			bound := &boundTraversalState{
				boundLayout:    l,
				section:        sec,
				currentAddress: address,
			}

			// apply assemble function
			assembled, err := assembleFunc(s, bound)
			if err != nil {
				return err
			}

			for _, a := range assembled {
				sec.AssembledStatements[j] = a
				j++
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

func (b *boundTraversalState) AddressFor(label *lexer.Token) (uint32, bool) {
	return b.boundLayout.LocateLabel(label.Value)
}

func (b *boundTraversalState) OffsetFor(label *lexer.Token) (uint16, error) {
	labelAddr, found := b.AddressFor(label)
	if !found {
		return 0, &asmerr.LabelNotFoundError{Label: label}
	}

	computed := int64(labelAddr) - int64(b.currentAddress)
	if computed > math.MaxUint16 || computed < 0 {
		return 0, &asmerr.BadComputedAddressError{
			Label:    label,
			Computed: computed,
			Signed:   false,
		}
	}

	return uint16(computed), nil
}

func (b *boundTraversalState) SignedOffsetFor(label *lexer.Token) (int16, error) {
	labelAddr, found := b.AddressFor(label)
	if !found {
		return 0, &asmerr.LabelNotFoundError{Label: label}
	}

	computed := int64(labelAddr) - int64(b.currentAddress)
	if computed > math.MaxInt16 || computed < math.MinInt16 {
		return 0, &asmerr.BadComputedAddressError{
			Label:    label,
			Computed: computed,
			Signed:   true,
		}
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
