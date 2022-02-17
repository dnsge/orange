package main

var instructionTokens = []*Instruction{
	// Identifiers
	{"REGISTER", "identifier", OneOf(`r[0-9]`, `r1[0-5]`), NoSlice},

	// Values
	{"BASE_8_IMM", "imm", Only(`#0o(-?[0-7]+)`), NoSlice},
	{"BASE_10_IMM", "imm", Only(`#(0|-?[1-9][0-9]*)`), NoSlice},
	{"BASE_16_IMM", "imm", Only(`#0x(-?[0-9A-Fa-f]+)`), NoSlice},
	{"STRING", NoCategory, OneOf(`"[^"]*"`, "`[^`]*`"), Slice(1, 1)},

	// Labels
	{"LABEL_DECLARATION", "directive", Only(`\$[a-zA-Z][a-zA-Z0-9]*:`), Slice(1, 1)},
	{"LABEL", NoCategory, Only(`\$[a-zA-Z][a-zA-Z0-9]*`), Slice(1, 0)},

	// Directives
	{"SECTION", "directive", Only(`\.section`), NoSlice},
	{"FILL_STATEMENT", "directive", Only(`\.fill`), NoSlice},
	{"STRING_STATEMENT", "directive", Only(`\.string`), NoSlice},

	// Whitespace + Formatting
	{"COMMA", NoCategory, Only(`,`), NoSlice},
	{"LBRACKET", NoCategory, Only(`\[`), NoSlice},
	{"RBRACKET", NoCategory, Only(`\]`), NoSlice},
	{"COMMENT", NoCategory, Only(`;[^\n]*`), NoSlice},
	{"LINE_END", NoCategory, Only(`\n`), NoSlice},

	// Operations
	{"ADD", OpCategory, DefaultPattern, NoSlice},
	{"ADDI", OpCategory, DefaultPattern, NoSlice},
	{"SUB", OpCategory, DefaultPattern, NoSlice},
	{"SUBI", OpCategory, DefaultPattern, NoSlice},
	{"AND", OpCategory, DefaultPattern, NoSlice},
	{"OR", OpCategory, DefaultPattern, NoSlice},
	{"XOR", OpCategory, DefaultPattern, NoSlice},
	{"LSL", OpCategory, DefaultPattern, NoSlice},
	{"LSR", OpCategory, DefaultPattern, NoSlice},
	{"CMP", OpCategory, DefaultPattern, NoSlice},
	{"CMPI", OpCategory, DefaultPattern, NoSlice},
	{"LDREG", OpCategory, DefaultPattern, NoSlice},
	{"LDWORD", OpCategory, DefaultPattern, NoSlice},
	{"LDHWRD", OpCategory, DefaultPattern, NoSlice},
	{"LDBYTE", OpCategory, DefaultPattern, NoSlice},
	{"STREG", OpCategory, DefaultPattern, NoSlice},
	{"STWORD", OpCategory, DefaultPattern, NoSlice},
	{"STHWRD", OpCategory, DefaultPattern, NoSlice},
	{"STBYTE", OpCategory, DefaultPattern, NoSlice},
	{"ADR", OpCategory, DefaultPattern, NoSlice},
	{"MOVZ", OpCategory, DefaultPattern, NoSlice},
	{"MOVK", OpCategory, DefaultPattern, NoSlice},
	{"B", OpCategory, DefaultPattern, NoSlice},
	{"BREG", OpCategory, DefaultPattern, NoSlice},
	{"B.EQ", OpCategory, DefaultPattern, NoSlice},
	{"B.NEQ", OpCategory, DefaultPattern, NoSlice},
	{"B.LT", OpCategory, DefaultPattern, NoSlice},
	{"B.LE", OpCategory, DefaultPattern, NoSlice},
	{"B.GT", OpCategory, DefaultPattern, NoSlice},
	{"B.GE", OpCategory, DefaultPattern, NoSlice},
	{"BL", OpCategory, DefaultPattern, NoSlice},
	{"PUSH", OpCategory, DefaultPattern, NoSlice},
	{"POP", OpCategory, DefaultPattern, NoSlice},
	{"SYSCALL", OpCategory, DefaultPattern, NoSlice},
	{"HALT", OpCategory, DefaultPattern, NoSlice},
	{"NOOP", OpCategory, DefaultPattern, NoSlice},

	// Generic identifier (last match)
	{"IDENTIFIER", "identifier", Only(`[a-zA-Z][a-zA-Z0-9]*`), NoSlice},
}
