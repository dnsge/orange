package asm

import (
	"fmt"
	"github.com/dnsge/orange/internal/asm/lexer"
	"github.com/dnsge/orange/internal/asm/parser"
)

func (a *assemblyContext) processLabelDeclarations() error {
	var address uint32 = 0
	for _, s := range a.statements {
		if s.Kind == parser.DirectiveStatement {
			directiveToken := s.Body[0]
			if directiveToken.Kind == lexer.LABEL_DECLARATION {
				labelName := directiveToken.Value[1:] // trim leading period
				if _, ok := a.labels[labelName]; ok {
					// duplicate label found
					return fmt.Errorf("duplicate label %q at %d:%d", labelName, directiveToken.Row, directiveToken.Column)
				}
				a.labels[labelName] = address
			}
		} else if s.Kind == parser.InstructionStatement {
			address++
		}
	}
	return nil
}
