package linker

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/linker/objfile"
	"io"
)

type AssembledSection struct {
	Name    string
	Size    int
	RawData []arch.Instruction

	// The absolute address where this section begins
	absoluteOffset int
}

type InputObjectFile struct {
	Sections        []*AssembledSection
	SymbolTable     []*objfile.SymbolTableEntry
	RelocationTable []*objfile.RelocationTableEntry
}

func readObjectFile(inputFile io.Reader) (*InputObjectFile, error) {
	of := &InputObjectFile{
		Sections:        nil,
		SymbolTable:     nil,
		RelocationTable: nil,
	}

	// scan sizes for slices
	var sectionCount, symbolCount, relocationCount int
	_, err := fmt.Fscanf(inputFile, "%d %d %d\n", &sectionCount, &symbolCount, &relocationCount)
	if err != nil {
		return nil, err
	}

	of.Sections = make([]*AssembledSection, sectionCount)
	of.SymbolTable = make([]*objfile.SymbolTableEntry, symbolCount)
	of.RelocationTable = make([]*objfile.RelocationTableEntry, relocationCount)

	// scan section definitions
	for i := range of.Sections {
		sec := new(AssembledSection)
		_, err = fmt.Fscanf(inputFile, "%s %d\n", &sec.Name, &sec.Size)
		if err != nil {
			return nil, err
		}

		// number of instructions is size divided by 4 bytes per instruction
		sec.RawData = make([]arch.Instruction, sec.Size/4)
		of.Sections[i] = sec
	}

	// scan symbol table entries
	for i := range of.SymbolTable {
		entry := new(objfile.SymbolTableEntry)
		err = entry.UnmarshalFrom(inputFile)
		if err != nil {
			return nil, err
		}
		of.SymbolTable[i] = entry
	}

	// scan relocation table entries
	for i := range of.RelocationTable {
		entry := new(objfile.RelocationTableEntry)
		err = entry.UnmarshalFrom(inputFile)
		if err != nil {
			return nil, err
		}
		of.RelocationTable[i] = entry
	}

	for _, sec := range of.Sections {
		for n := 0; n < len(sec.RawData); n++ {
			err = binary.Read(inputFile, binary.LittleEndian, &sec.RawData[n])
			if err != nil {
				return nil, err
			}
		}
	}

	return of, nil
}

func (i *InputObjectFile) GetSymbolAbsoluteAddress(symbolName string) (int, error) {
	entry, ok := i.getSymbolEntryByName(symbolName)
	if !ok {
		return 0, fmt.Errorf("symbol %q not found in object file", symbolName)
	}

	section, ok := i.getSectionByName(entry.SectionName)
	if !ok {
		return 0, fmt.Errorf("symbol %q is in non-existant section %q", symbolName, entry.SectionName)
	}

	totalOffset := section.absoluteOffset + entry.SectionOffset
	return totalOffset, nil
}

func (i *InputObjectFile) GetSectionOffsetAbsoluteAddress(sectionName string, sectionOffset int) (int, error) {
	section, ok := i.getSectionByName(sectionName)
	if !ok {
		return 0, fmt.Errorf("section %q does not exist", sectionName)
	}

	offset := section.absoluteOffset + sectionOffset
	return offset, nil
}

func (i *InputObjectFile) getSymbolEntryByName(symbolName string) (*objfile.SymbolTableEntry, bool) {
	for _, entry := range i.SymbolTable {
		if entry.LabelName == symbolName {
			return entry, true
		}
	}
	return nil, false
}

func (i *InputObjectFile) getSectionByName(section string) (*AssembledSection, bool) {
	for _, s := range i.Sections {
		if s.Name == section {
			return s, true
		}
	}
	return nil, false
}
