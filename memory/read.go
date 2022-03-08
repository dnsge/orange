package memory

import (
	"bytes"
	"io"
)

// LoadFromReader loads into memory from an io.Reader, allocating a new block at startAddress.
// The memory must not have been already allocated to another block.
func (m *Memory) LoadFromReader(startAddress uint32, reader io.Reader) (int, error) {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return 0, err
	}

	b := allocateBlock(startAddress, uint32(buf.Len()))
	copy(b.data, buf.Bytes())
	m.Blocks = append(m.Blocks, b)
	return buf.Len(), nil
}
