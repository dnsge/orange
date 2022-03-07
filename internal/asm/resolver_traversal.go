package asm

import (
	"errors"
	"github.com/dnsge/orange/internal/asm/asmerr"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
	"github.com/dnsge/orange/internal/linker/objfile"
)

// resolverTraversalState implements and wraps a TraversalState to provide
// label resolution for undefined labels in the case of generating an
// object file. Additionally, it automatically adds requested addresses
// to the relocation table of the object file.
type resolverTraversalState struct {
	layout     *Layout
	state      TraversalState
	statement  *parser.Statement
	objectFile *ObjectFile
}

// addUnresolvedLabel adds the given label to the object file's symbol table,
// marking it as unresolved. The section offset is temporarily resolved to 0
// and will be updated at link time.
func (r *resolverTraversalState) addUnresolvedLabel(label *lexer.Token) {
	// add unresolved entry to symbol table
	r.objectFile.optionallyAddSymbol(label.Value, &objfile.SymbolTableEntry{
		LabelName:     label.Value,
		SectionName:   "-",
		SectionOffset: 0,
		Resolved:      false,
	})

	// unresolved label value must be relocated at link time
	r.addCurrentToRelocationTable(label)
}

// addCurrentToRelocationTable marks the current statement (at address) as
// needing relocation. This is necessary for absolute address values and
// for relative addresses between sections, as sections can get reordered
// and resized at link time.
func (r *resolverTraversalState) addCurrentToRelocationTable(label *lexer.Token) {
	r.objectFile.RelocationTable = append(r.objectFile.RelocationTable, &objfile.RelocationTableEntry{
		LabelName:      label.Value,
		SectionName:    r.Section().Name,
		SectionOffset:  r.Address(),
		StatementToken: r.statement.Body[0].Kind,
	})
}

func (r *resolverTraversalState) AddressFor(label *lexer.Token) (uint32, bool) {
	res, ok := r.state.AddressFor(label)
	if !ok {
		if isPrivateLabel(label) {
			// we must locate private labels, so return not found
			return 0, false
		}
		r.addUnresolvedLabel(label)
		return 0, true
	}

	// add to relocation table regardless because absolute addresses are very
	// likely to change at link time
	r.addCurrentToRelocationTable(label)
	return res, true
}

func (r *resolverTraversalState) OffsetFor(label *lexer.Token) (uint16, error) {
	res, err := r.state.OffsetFor(label)

	// Check if LabelNotFoundError and add unresolved if so
	var nfErr *asmerr.LabelNotFoundError
	if errors.As(err, &nfErr) {
		if isPrivateLabel(label) {
			// we must locate private labels, so return an error
			return 0, nfErr
		}
		r.addUnresolvedLabel(label)
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	// we can ignore the ok result because we have already cleared that the variable exists
	labelSection, _ := r.layout.LocateLabelSection(label.Value)
	if labelSection != r.state.Section() {
		// different sections, so add to relocation table
		r.addCurrentToRelocationTable(label)
	}

	return res, nil
}

func (r *resolverTraversalState) SignedOffsetFor(label *lexer.Token) (int16, error) {
	res, err := r.state.SignedOffsetFor(label)

	// Check if LabelNotFoundError and add unresolved if so
	var nfErr *asmerr.LabelNotFoundError
	if errors.As(err, &nfErr) {
		if isPrivateLabel(label) {
			// we must locate private labels, so return an error
			return 0, nfErr
		}
		r.addUnresolvedLabel(label)
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	// we can ignore the ok result because we have already cleared that the variable exists
	labelSection, _ := r.layout.LocateLabelSection(label.Value)
	if labelSection != r.state.Section() {
		// different sections, so add to relocation table
		r.addCurrentToRelocationTable(label)
	}

	return res, nil
}

func (r *resolverTraversalState) Section() *Section {
	return r.state.Section()
}

func (r *resolverTraversalState) Address() int {
	return r.state.Address()
}

func (r *resolverTraversalState) AdvanceAddress(amount int) {
	r.state.AdvanceAddress(amount)
}

// isPrivateLabel returns whether the given label is only visible within
// the current file, denoted by an underscore preceding the name.
func isPrivateLabel(label *lexer.Token) bool {
	return label.Value[0] == '_'
}
