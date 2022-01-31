package memory

import "fmt"

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
		return joinByteSlice(b.data[dataStart : dataStart+2])
	case 32: // word read
		return joinByteSlice(b.data[dataStart : dataStart+4])
	case 64: // double word (register) read
		return joinByteSlice(b.data[dataStart : dataStart+8])
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
	split := splitUint64(data, bytes)
	copy(b.data[dataStart:dataEnd], split)
}

func joinByteSlice(slice []byte) uint64 {
	if len(slice) > 8 {
		panic("can't join more than 8 bytes into uint64")
	} else if len(slice) == 0 {
		return 0
	}

	var res uint64 = 0
	for i := len(slice) - 1; i >= 0; i-- {
		res |= uint64(slice[i])
		if i != 0 {
			res <<= 8
		}
	}
	return res
}

func splitUint64(num uint64, bytes uint32) (res []byte) {
	res = make([]byte, bytes)
	for i := uint32(0); i < bytes; i++ {
		res[i] = byte(num & 0xFF)
		num >>= 8
	}
	return
}
