package linker

import (
	"encoding/binary"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"io"
)

type linkContext struct {
	Symbols      map[string]*InputObjectFile
	Instructions []arch.Instruction
}

func Link(inputFiles []io.Reader, outputFile io.Writer) error {
	objectFiles := make([]*InputObjectFile, len(inputFiles))
	for i, f := range inputFiles {
		obj, err := readObjectFile(f)
		if err != nil {
			return err
		}
		objectFiles[i] = obj
	}

	collectedSymbols, err := collectSymbolTableEntries(objectFiles)
	if err != nil {
		return err
	}

	collectedSections, sectionOrder := collectAssembledSections(objectFiles)
	instructions := layoutAllSections(collectedSections, sectionOrder)

	linkCtx := &linkContext{
		Symbols:      collectedSymbols,
		Instructions: instructions,
	}

	err = linkCtx.relocateAll(objectFiles)
	if err != nil {
		return err
	}

	err = linkCtx.writeInstructions(outputFile)
	return err
}

// collectSymbolTableEntries groups all the symbols from the input object
// files' symbol tables. Unresolved symbols (e.g. requested by a file but
// never defined) or duplicate symbols (e.g. two matching labels in different
// files) will return an error.
func collectSymbolTableEntries(objectFiles []*InputObjectFile) (map[string]*InputObjectFile, error) {
	res := make(map[string]*InputObjectFile)
	for _, objFile := range objectFiles {
		for _, symbol := range objFile.SymbolTable {
			// If the symbol is already in res and was resolved in a different file
			if existing, ok := res[symbol.LabelName]; ok && existing != nil {
				if symbol.Resolved {
					// duplicate symbol entry
					return nil, fmt.Errorf("duplicate symbol %q", symbol.LabelName)
				} else {
					// Symbol is not resolved in this file, but it was already resolved
					// elsewhere. Just continue.
					continue
				}
			}

			if symbol.Resolved {
				// Mark as resolved
				res[symbol.LabelName] = objFile
			} else {
				// Mark as still unresolved
				res[symbol.LabelName] = nil
			}
		}
	}

	// Check for any nil entries which signifies an undefined symbol
	for name, objFile := range res {
		if objFile == nil {
			return nil, fmt.Errorf("undefined symbol %q", name)
		}
	}

	return res, nil
}

// collectAssembledSections groups all the AssembledSections specified by the
// input object files, returning a map of every section name to the array of
// section instances and the order that the section names appeared.
func collectAssembledSections(objectFiles []*InputObjectFile) (map[string][]*AssembledSection, []string) {
	// use map and order slice to imitate ordered map behavior
	sections := make(map[string][]*AssembledSection)
	var order []string
	for _, objFile := range objectFiles {
		for _, section := range objFile.Sections {
			sectionGroup := sections[section.Name]
			if sectionGroup == nil {
				// first of this section name, add to order
				order = append(order, section.Name)
			}
			sections[section.Name] = append(sectionGroup, section)
		}
	}

	return sections, order
}

func layoutAllSections(sections map[string][]*AssembledSection, sectionOrder []string) (res []arch.Instruction) {
	address := 0
	for _, sectionName := range sectionOrder {
		assembledSections := sections[sectionName]
		// iterate over each defined section that shares the same name
		for _, section := range assembledSections {
			res = append(res, section.RawData...)
			section.absoluteOffset = address
			address += section.Size
		}
	}
	return res
}

// relocateAll performs the relocation for every object file given. Each symbol
// specified in the InputObjectFile's relocation table is found and the
// corresponding instruction is then updated.
func (l *linkContext) relocateAll(objectFiles []*InputObjectFile) error {
	for _, objFile := range objectFiles {
		for _, relocation := range objFile.RelocationTable {
			symbolFile, ok := l.Symbols[relocation.LabelName]
			if !ok {
				return fmt.Errorf("relocation table contains undefined symbol %q", relocation.LabelName)
			}

			// Find the address of the symbol to relocate for
			symbolAddress, err := symbolFile.GetSymbolAbsoluteAddress(relocation.LabelName)
			if err != nil {
				return err
			}

			// Compute the absolute address of the instruction being relocated
			relocationAddress, err := objFile.GetSectionOffsetAbsoluteAddress(relocation.SectionName, relocation.SectionOffset)
			if err != nil {
				return err
			}

			// Actually modify the instruction from linkContext.Instructions
			err = performRelocation(l.Instructions, symbolAddress, relocationAddress, relocation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// writeInstructions writes each instruction to the output writer
func (l *linkContext) writeInstructions(writer io.Writer) error {
	for i := range l.Instructions {
		if err := binary.Write(writer, arch.ByteOrder, l.Instructions[i]); err != nil {
			return err
		}
	}
	return nil
}
