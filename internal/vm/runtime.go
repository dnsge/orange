package vm

import (
	"errors"
	"fmt"
	"github.com/dnsge/orange/internal/arch"
	"log"
)

const (
	syscallRead  = 0
	syscallWrite = 1
)

var (
	ErrInvalidFileDescriptor = fmt.Errorf("invalid file descriptor")
)

func (v *VirtualMachine) executeSyscall() {
	syscallNumber := v.registers.Get(arch.SyscallRegister)
	if !v.quiet {
		log.Printf("Executing syscall number %d\n", syscallNumber)
	}
	switch syscallNumber {
	case syscallRead:
		v.syscallRead()
	case syscallWrite:
		v.syscallWrite()
	default:
		panic(fmt.Sprintf("invalid syscall number %d", syscallNumber))
	}
}

// Available Syscalls:
// - read(int fileD, char *buf, size_t bytes)
// - write(int fileD, const char *buf, size_t bytes)

// syscallRead performs a generic read syscall
//
// Semantics: write(int fileD, const void *buf, size_t nBytes)
// fileD:  register 1
// buf:    register 2
// nBytes: register 3
//
// Reads the first nBytes from the file into buf.
func (v *VirtualMachine) syscallRead() {
	fileD := int(v.registers.Get(1))
	bufPtr := uint32(v.registers.Get(2))
	nBytes := uint32(v.registers.Get(3))

	err := v.syscallReadExecute(fileD, bufPtr, nBytes)
	if err != nil {
		if errors.Is(err, ErrInvalidFileDescriptor) {
			v.setSyscallError(arch.ENO_BadFileDescriptor)
		} else {
			v.setSyscallError(arch.ENO_IO)
		}
		log.Printf("error: %v\n", err)
	}
}

func (v *VirtualMachine) syscallReadExecute(fileD int, bufPtr uint32, nBytes uint32) error {
	file, ok := v.fds[fileD]
	if !ok {
		return ErrInvalidFileDescriptor
	}

	buf := make([]byte, nBytes)
	_, err := file.Read(buf)
	if err != nil {
		return err // todo: EOF
	}

	for i := uint32(0); i < nBytes; i++ {
		v.memory.Write(bufPtr+i, 8, uint64(buf[i]))
	}

	return nil
}

// syscallWrite performs a generic write syscall
//
// Semantics: write(int fileD, const void *buf, size_t nBytes)
// fileD:  register 1
// buf:    register 2
// nBytes: register 3
//
// Writes the first nBytes bytes of buf to the file.
func (v *VirtualMachine) syscallWrite() {
	fileD := int(v.registers.Get(1))
	bufPtr := uint32(v.registers.Get(2))
	nBytes := uint32(v.registers.Get(3))

	err := v.syscallWriteExecute(fileD, bufPtr, nBytes)
	if err != nil {
		if errors.Is(err, ErrInvalidFileDescriptor) {
			v.setSyscallError(arch.ENO_BadFileDescriptor)
		} else {
			v.setSyscallError(arch.ENO_IO)
		}
		log.Printf("error: %v\n", err)
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

func (v *VirtualMachine) setSyscallError(eno uint64) {
	v.registers.Set(arch.SyscallErrorRegister, eno)
}
