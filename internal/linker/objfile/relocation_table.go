package objfile

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"io"
)

type RelocationTableEntry struct {
	LabelName      string
	SectionName    string
	SectionOffset  int
	StatementToken lexer.TokenKind
}

func (r *RelocationTableEntry) String() string {
	return fmt.Sprintf("[%s@%s : offset=%d, kind=%s]", r.LabelName, r.SectionName, r.SectionOffset, lexer.DescribeTokenKind(r.StatementToken))
}

func (r *RelocationTableEntry) MarshalTo(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "%s %s %d %d\n", r.LabelName, r.SectionName, r.SectionOffset, r.StatementToken)
	return err
}

func (r *RelocationTableEntry) UnmarshalFrom(reader io.Reader) error {
	_, err := fmt.Fscanf(reader, "%s %s %d %d\n", &r.LabelName, &r.SectionName, &r.SectionOffset, &r.StatementToken)
	return err
}
