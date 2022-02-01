package memory

import (
	"encoding/binary"
	"fmt"
)

var (
	byteOrder = binary.LittleEndian
)

// block describes a virtual memory space containing [startAddress, endAddress)
type block struct {
	startAddress uint32
	endAddress   uint32
	data         []byte
}

func allocateBlock(startAddress uint32, size uint32) *block {
	return &block{
		startAddress: startAddress,
		endAddress:   startAddress + size - 1,
		data:         make([]byte, size),
	}
}

func (b *block) Contains(address uint32) bool {
	return address >= b.startAddress && address < b.endAddress
}

func (b *block) Read(address uint32, size uint32) uint64 {
	dataStart := address - b.startAddress
	switch size {
	case 8: // byte read
		return uint64(b.data[dataStart])
	case 16: // half-word read
		return uint64(byteOrder.Uint16(b.data[dataStart : dataStart+2]))
	case 32: // word read
		return uint64(byteOrder.Uint32(b.data[dataStart : dataStart+4]))
	case 64: // double word (register) read
		return byteOrder.Uint64(b.data[dataStart : dataStart+8])
	default:
		panic(fmt.Sprintf("invalid read size of %d", size))
	}
}

func (b *block) Write(address uint32, size uint32, data uint64) {
	if size%8 != 0 || size == 0 {
		panic(fmt.Sprintf("invalid write size of %d", size))
	}
	bytes := size / 8
	dataStart := address - b.startAddress
	dataEnd := dataStart + bytes
	split := make([]byte, 8)
	byteOrder.PutUint64(split, data)
	copy(b.data[dataStart:dataEnd], split)
}
