package arch

const (
	opcodeOffset = 24

	a_destRegOffset = 20
	a_aRegOffset    = 16
	a_bRegOffset    = 12
	a_immOffset     = 0

	m_aRegOffset = 20
	m_bRegOffset = 16
	m_immOffset  = 0

	e_destRegOffset = 20
	e_immOffset     = 0

	b_aRegOffset = 20
	b_immOffset  = 0

	r_aRegOffset = 20

	regValMask = 0xF
	immValMask = 0xFFFF
)

type Instruction = uint32

func GetOpcode(instruction Instruction) Opcode {
	opcodeVal := (instruction >> opcodeOffset) & 0x3F
	return Opcode(opcodeVal)
}

type ATypeInstruction struct {
	Opcode  Opcode
	RegDest RegisterValue
	RegA    RegisterValue
	RegB    RegisterValue
}

func DecodeATypeInstruction(instruction Instruction, opcode Opcode) ATypeInstruction {
	return ATypeInstruction{
		Opcode:  opcode,
		RegDest: extractRegister(instruction, a_destRegOffset),
		RegA:    extractRegister(instruction, a_aRegOffset),
		RegB:    extractRegister(instruction, a_bRegOffset),
	}
}

func EncodeATypeInstruction(instruction ATypeInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= uint32(instruction.RegDest) << a_destRegOffset
	encoded |= uint32(instruction.RegA) << a_aRegOffset
	encoded |= uint32(instruction.RegB) << a_bRegOffset
	return encoded
}

type ATypeImmInstruction struct {
	Opcode    Opcode
	RegDest   RegisterValue
	RegA      RegisterValue
	Immediate uint16
}

func DecodeATypeImmInstruction(instruction Instruction, opcode Opcode) ATypeImmInstruction {
	return ATypeImmInstruction{
		Opcode:    opcode,
		RegDest:   extractRegister(instruction, a_destRegOffset),
		RegA:      extractRegister(instruction, a_aRegOffset),
		Immediate: extractUnsignedImm(instruction, a_immOffset),
	}
}

func EncodeATypeImmInstruction(instruction ATypeImmInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= uint32(instruction.RegDest) << a_destRegOffset
	encoded |= uint32(instruction.RegA) << a_aRegOffset
	encoded |= uint32(instruction.Immediate) << a_immOffset
	return encoded
}

type MTypeInstruction struct {
	Opcode    Opcode
	RegA      RegisterValue
	RegB      RegisterValue
	Immediate int16
}

func DecodeMTypeInstruction(instruction Instruction, opcode Opcode) MTypeInstruction {
	return MTypeInstruction{
		Opcode:    opcode,
		RegA:      extractRegister(instruction, m_aRegOffset),
		RegB:      extractRegister(instruction, m_bRegOffset),
		Immediate: extractSignedImm(instruction, m_immOffset),
	}
}

func EncodeMTypeInstruction(instruction MTypeInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= uint32(instruction.RegA) << m_aRegOffset
	encoded |= uint32(instruction.RegB) << m_bRegOffset
	encoded |= (uint32(instruction.Immediate) & 0xFFFF) << m_immOffset
	return encoded
}

type ETypeInstruction struct {
	Opcode    Opcode
	RegDest   RegisterValue
	Immediate uint16
}

func DecodeETypeInstruction(instruction Instruction, opcode Opcode) ETypeInstruction {
	return ETypeInstruction{
		Opcode:    opcode,
		RegDest:   extractRegister(instruction, e_destRegOffset),
		Immediate: extractUnsignedImm(instruction, e_immOffset),
	}
}

func EncodeETypeInstruction(instruction ETypeInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= uint32(instruction.RegDest) << e_destRegOffset
	encoded |= uint32(instruction.Immediate) << e_immOffset
	return encoded
}

type BTypeInstruction struct {
	Opcode Opcode
	RegA   RegisterValue
}

func DecodeBTypeInstruction(instruction Instruction, opcode Opcode) BTypeInstruction {
	return BTypeInstruction{
		Opcode: opcode,
		RegA:   extractRegister(instruction, b_aRegOffset),
	}
}

func EncodeBTypeInstruction(instruction BTypeInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= uint32(instruction.RegA) << b_aRegOffset
	return encoded
}

type BTypeImmInstruction struct {
	Opcode Opcode
	Offset int16
}

func DecodeBTypeImmInstruction(instruction Instruction, opcode Opcode) BTypeImmInstruction {
	return BTypeImmInstruction{
		Opcode: opcode,
		Offset: extractSignedImm(instruction, b_immOffset),
	}
}

func EncodeBTypeImmInstruction(instruction BTypeImmInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= (uint32(instruction.Offset) & 0xFFFF) << b_immOffset
	return encoded
}

type RTypeInstruction struct {
	Opcode Opcode
	RegA   RegisterValue
}

func DecodeRTypeInstruction(instruction Instruction, opcode Opcode) RTypeInstruction {
	return RTypeInstruction{
		Opcode: opcode,
		RegA:   extractRegister(instruction, r_aRegOffset),
	}
}

func EncodeRTypeInstruction(instruction RTypeInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	encoded |= uint32(instruction.RegA) << r_aRegOffset
	return encoded
}

type OTypeInstruction struct {
	Opcode Opcode
}

func DecodeOTypeInstruction(instruction Instruction, opcode Opcode) OTypeInstruction {
	return OTypeInstruction{
		Opcode: opcode,
	}
}

func EncodeOTypeInstruction(instruction OTypeInstruction) Instruction {
	var encoded uint32 = 0
	encoded |= uint32(instruction.Opcode) << opcodeOffset
	return encoded
}

func extractRegister(instruction Instruction, offset uint8) RegisterValue {
	return RegisterValue((instruction >> offset) & regValMask)
}

func extractUnsignedImm(instruction Instruction, offset uint8) uint16 {
	return uint16((instruction >> offset) & immValMask)
}

func extractSignedImm(instruction Instruction, offset uint8) int16 {
	return int16((instruction >> offset) & immValMask)
}
