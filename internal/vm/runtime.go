package vm

import (
	"fmt"
	"github.com/dnsge/orange/internal/arch"
)

const (
	syscallWrite = 1
)

var (
	ErrInvalidFileDescriptor = fmt.Errorf("invalid file descriptor")
)

func (v *VirtualMachine) executeSyscall() {
	syscallNumber := v.registers.Get(arch.SyscallRegister)
	switch syscallNumber {
	case syscallWrite:
		v.syscallWrite()
	default:
		panic(fmt.Sprintf("invalid syscall number %d", syscallNumber))
	}
}

// Available Syscalls:
// - write(int fileD, const void *buf, size_t bytes)

// syscallWrite performs a generic write syscall
//
// Semantics: write(int fileD, const void *buf, size_t nBytes)
// fileD:  register 1
// buf:    register 2
// nBytes: register 3
//
// Writes the first nBytes bytes of buf to the file.
// In true
func (v *VirtualMachine) syscallWrite() {
	fileD := int(v.registers.Get(1))
	bufPtr := uint32(v.registers.Get(2))
	nBytes := uint32(v.registers.Get(3))

	err := v.syscallWriteExecute(fileD, bufPtr, nBytes)
	if err != nil {
		panic(err)
	}
}

func (v *VirtualMachine) syscallWriteExecute(fileD int, bufPtr uint32, nBytes uint32) error {
	file, ok := v.fds[fileD]
	if !ok {
		return ErrInvalidFileDescriptor
	}

	for i := uint32(0); i < nBytes; i++ {
		data := v.memory.Read(bufPtr, 8) // read single byte
		singleByte := byte(data)
		_, err := file.Write([]byte{singleByte})
		if err != nil {
			return fmt.Errorf("syscall write: %w", err)
		}
		bufPtr += 1
	}

	return nil
}
