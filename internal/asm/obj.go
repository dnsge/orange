package asm

type SymbolTableEntry struct {
	LabelName     string
	SectionName   string
	SectionOffset int
}

type RelocationTableEntry struct {
	SectionName   string
	SectionOffset int
}

type ObjectFile struct {
	Sections    []Section
	SymbolTable []SymbolTableEntry
}
