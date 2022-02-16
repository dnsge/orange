package vm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"github.com/dnsge/orange/internal/memory"
)

type VirtualMachine struct {
	programCounter uint32
	registers      registerFile
	flags          flags
	memory         memory.Addressable
	halted         bool
}

func (v *VirtualMachine) Memory() memory.Addressable {
	return v.memory
}

func (v *VirtualMachine) Halt() {
	v.halted = true
}

func (v *VirtualMachine) Halted() bool {
	return v.halted
}

func NewVirtualMachine(mem memory.Addressable) *VirtualMachine {
	return &VirtualMachine{
		programCounter: 0,
		registers:      initRegisterFile(),
		flags: flags{
			Negative: false,
			Zero:     false,
			Carry:    false,
		},
		memory: mem,
		halted: false,
	}
}

func (v *VirtualMachine) InitStack(stackStartAddress uint64) {
	v.registers.Set(arch.StackRegister, stackStartAddress)
}

func (v *VirtualMachine) ExecuteInstruction() {
	if v.halted {
		return
	}

	i := v.fetchNextInstruction()
	v.executeInstruction(i)
}

func (v *VirtualMachine) PrintState() {
	fmt.Printf("Registers: %v\n", v.registers)
	fmt.Printf("PC: 0x%08x\n", v.programCounter)
}

func (v *VirtualMachine) setFlags(res, carry uint64) {
	v.flags.Zero = res == 0
	v.flags.Negative = res&(0b1000<<60) > 0 // check last bit for signed-ness
	v.flags.Carry = carry == 1
}
