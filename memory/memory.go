package memory

import "fmt"

type Addressable interface {
	Read(address uint32, size uint32) uint64
	Write(address uint32, size uint32, data uint64)
}

type Memory struct {
	Blocks []*block
}

func New() *Memory {
	return &Memory{}
}

func (m *Memory) Alloc(startAddress uint32, size uint32) {
	m.Blocks = append(m.Blocks, allocateBlock(startAddress, size))
}

func (m *Memory) Read(address uint32, size uint32) uint64 {
	for _, b := range m.Blocks {
		if b.Contains(address) {
			return b.Read(address, size)
		}
	}
	panic(fmt.Sprintf("invalid read at address 0x%08x", address))
}

func (m *Memory) Write(address uint32, size uint32, data uint64) {
	for _, b := range m.Blocks {
		if b.Contains(address) {
			b.Write(address, size, data)
			return
		}
	}
	panic(fmt.Sprintf("invalid write at address 0x%08x", address))
}
