package memory

import "encoding/binary"

const (
	chunkSize  = 8
	chunkMask  = 7
	chunkShift = 3
)

func Load(offset ExtismPointer, buf []byte) {
	length := len(buf)
	chunkCount := length >> chunkShift

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << chunkShift
		binary.LittleEndian.PutUint64(buf[i:i+chunkSize], ExtismLoadU64(offset+ExtismPointer(i)))
	}

	remainder := length & chunkMask
	remainderOffset := chunkCount << chunkShift
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		buf[index] = ExtismLoadU8(offset + ExtismPointer(index))
	}
}

func Store(offset ExtismPointer, buf []byte) {
	length := len(buf)
	chunkCount := length >> chunkShift

	for chunkIdx := 0; chunkIdx < chunkCount; chunkIdx++ {
		i := chunkIdx << chunkShift
		x := binary.LittleEndian.Uint64(buf[i : i+chunkSize])
		ExtismStoreU64(offset+ExtismPointer(i), x)
	}

	remainder := length & chunkMask
	remainderOffset := chunkCount << chunkShift
	for index := remainderOffset; index < (remainder + remainderOffset); index++ {
		ExtismStoreU8(offset+ExtismPointer(index), buf[index])
	}
}

func NewMemory(offset ExtismPointer, length uint64) Memory {
	return Memory{
		offset: offset,
		length: length,
	}
}

// Memory represents memory allocated by (and shared with) the host.
type Memory struct {
	offset ExtismPointer
	length uint64
}

// Load copies the host memory block to the provided `buffer` byte slice.
func (m *Memory) Load(buffer []byte) {
	Load(m.offset, buffer)
}

// Store copies the `data` byte slice into host memory.
func (m *Memory) Store(data []byte) {
	Store(m.offset, data)
}

// Free frees the host memory block.
func (m *Memory) Free() {
	ExtismFree(m.offset)
}

// Length returns the number of bytes in the host memory block.
func (m *Memory) Length() uint64 {
	return m.length
}

// Offset returns the offset of the host memory block.
func (m *Memory) Offset() uint64 {
	return uint64(m.offset)
}

// ReadBytes returns the host memory block as a slice of bytes.
func (m *Memory) ReadBytes() []byte {
	buff := make([]byte, m.length)
	m.Load(buff)
	return buff
}
