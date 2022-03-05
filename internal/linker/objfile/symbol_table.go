package objfile

import (
	"fmt"
	"io"
)

type SymbolTableEntry struct {
	LabelName     string
	SectionName   string
	SectionOffset int
	Resolved      bool
}

func (s *SymbolTableEntry) String() string {
	return fmt.Sprintf("[%s@%s : offset=%d, resolved=%t]", s.LabelName, s.SectionName, s.SectionOffset, s.Resolved)
}

func (s *SymbolTableEntry) MarshalTo(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "%s %s %d %t\n", s.LabelName, s.SectionName, s.SectionOffset, s.Resolved)
	return err
}

func (s *SymbolTableEntry) UnmarshalFrom(reader io.Reader) error {
	_, err := fmt.Fscanf(reader, "%s %s %d %t\n", &s.LabelName, &s.SectionName, &s.SectionOffset, &s.Resolved)
	return err
}
