package asm

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/parser"
	"github.com/dnsge/orange/internal/linker/objfile"
	"io"
)

type ObjectFile struct {
	Sections        []*Section
	SymbolTable     []*objfile.SymbolTableEntry
	RelocationTable []*objfile.RelocationTableEntry

	symbolTableMap map[string]struct{}
}

func printObjectFile(of *ObjectFile) {
	for _, s := range of.SymbolTable {
		fmt.Printf("symbol %s\n", s)
	}

	for _, r := range of.RelocationTable {
		fmt.Printf("relocation %s\n", r)
	}
}

func CreateObjectFile(layout *Layout) (*ObjectFile, error) {
	of := &ObjectFile{
		Sections:        nil,
		SymbolTable:     nil,
		RelocationTable: nil,
		symbolTableMap:  make(map[string]struct{}),
	}

	// copy sections
	of.Sections = make([]*Section, len(layout.Sections))
	for i, sec := range layout.Sections {
		of.Sections[i] = sec
	}

	for labelName, statement := range layout.Labels {
		offset, section, ok := layout.LocateStatementWithinSection(statement)
		if !ok {
			panic("failed to locate statement in section") // todo: Proper error instead of panicking
		}

		of.optionallyAddSymbol(labelName, &objfile.SymbolTableEntry{
			LabelName:     labelName,
			SectionName:   section.Name,
			SectionOffset: offset,
			Resolved:      true,
		})
	}

	return of, nil
}

// optionallyAddSymbol adds the given objfile.SymbolTableEntry to the ObjectFile's
// array of symbols if the label name has not already been added.
func (o *ObjectFile) optionallyAddSymbol(labelName string, entry *objfile.SymbolTableEntry) {
	if _, ok := o.symbolTableMap[labelName]; ok {
		// already added this symbol to the symbol table
		return
	}

	o.SymbolTable = append(o.SymbolTable, entry)
	o.symbolTableMap[labelName] = struct{}{}
}

func (o *ObjectFile) AssembleStatement(layout *Layout) func(s *parser.Statement, state TraversalState) ([]arch.Instruction, error) {
	return func(s *parser.Statement, state TraversalState) ([]arch.Instruction, error) {
		// Construct resolverTraversalState with encapsulated layout
		resolver := &resolverTraversalState{
			layout:     layout,
			state:      state,
			statement:  s,
			objectFile: o,
		}

		// transparently use typical assembly function
		return AssembleStatement(s, resolver)
	}
}

// WriteToFile writes the ObjectFile to the given io.Writer completely.
//
// The file format is as follows:
//
// [# of sections] [# of symbol table entries] [# of relocation table entries]
// - for each section, [section name] [section size]
// - for each entry, [label name] [section name] [offset] [resolved]
// - for each entry, [label name] [section name] [offset] [instruction/directive id]
// [raw assembled instructions]
func (o *ObjectFile) WriteToFile(layout *Layout, outputFile io.Writer) (err error) {
	// Write # of sections, size of symbol table, size of relocation table
	_, err = fmt.Fprintf(outputFile, "%d %d %d\n", len(o.Sections), len(o.SymbolTable), len(o.RelocationTable))
	if err != nil {
		return
	}

	// Write section sizes
	for _, sec := range o.Sections {
		_, err = fmt.Fprintf(outputFile, "%s %d\n", sec.Name, sec.Size)
		if err != nil {
			return
		}
	}

	// Write symbol table entries
	for _, entry := range o.SymbolTable {
		err = entry.MarshalTo(outputFile)
		if err != nil {
			return
		}
	}

	// Write relocation table entries
	for _, entry := range o.RelocationTable {
		err = entry.MarshalTo(outputFile)
		if err != nil {
			return
		}
	}

	// write assembled statements
	err = layout.Traverse(func(section *Section) error {
		for _, a := range section.AssembledStatements {
			err = binary.Write(outputFile, arch.ByteOrder, a)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return
}
