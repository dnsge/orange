package main

import (
	"fmt"
	"strings"
)

const (
	NoCategory = ""
	OpCategory = "op"
)

type Instruction struct {
	TokenName     string
	TokenCategory string
	Pattern       Pattern
	Slice         TokenSlice
}

// EnumName returns the source code token name, like B_EQ
func (i *Instruction) EnumName() string {
	return strings.ReplaceAll(i.TokenName, ".", "_")
}

// CaptureFunc returns the proper token capture function based on the
// set TokenSlice
func (i *Instruction) CaptureFunc() string {
	if i.Slice == NoSlice {
		return fmt.Sprintf(`tokenOfKind(%s)`, i.EnumName())
	} else {
		return fmt.Sprintf(`tokenOfKindSliced(%s, %d, %d)`, i.EnumName(), i.Slice.Start, i.Slice.End)
	}
}

// TokenRegex returns a regex-form of the default token name.
// This is only used in the case of DefaultPattern
func (i *Instruction) TokenRegex() string {
	return strings.ReplaceAll(i.TokenName, ".", "\\.")
}

// fakeOps are operations that are translated away and thus need no place
// in the conversion to arch opcode entries.
var fakeOps = []string{
	"ADR",
	"CMP",
	"CMPI",
	"MOV",
}

func (i *Instruction) IsRealOp() bool {
	if i.TokenCategory != OpCategory {
		return false
	}

	for _, fake := range fakeOps {
		if fake == i.TokenName {
			return false
		}
	}
	return true
}
