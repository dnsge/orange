package asmerr

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/asm/lexer"
)

type InvalidArgumentCountError struct {
	Opcode   arch.Opcode
	Expected int
	Got      int
}

func (i *InvalidArgumentCountError) Error() string {
	return fmt.Sprintf("invalid argument count for %s: expected %d but got %d",
		i.Opcode, i.Expected, i.Got)
}

type BadComputedAddressError struct {
	Label    *lexer.Token
	Computed int64
	Signed   bool
}

func (b *BadComputedAddressError) Error() string {
	return fmt.Sprintf("cannot represent computed value for label %s as %s: computed value is %d",
		describeLocatedToken(b.Label), b.signText(), b.Computed)
}

func (b *BadComputedAddressError) signText() string {
	if b.Signed {
		return "signed 16-bit integer"
	} else {
		return "unsigned 16-bit integer"
	}
}
