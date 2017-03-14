package cli

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ShapeReader struct {
	reader           *bytes.Reader
	originalPosition int64
}

func NewShapeReader(reader *bytes.Reader) *ShapeReader {
	originalPosition, _ := reader.Seek(0, io.SeekCurrent)
	metaReader := &ShapeReader{
		reader:           reader,
		originalPosition: originalPosition,
	}
	return metaReader
}

// Seek implements the io.Seeker interface.
func (sr *ShapeReader) Seek(offset int64) error {
	_, err := sr.reader.Seek(sr.originalPosition+offset, io.SeekStart)
	return err
}

// ReadByte =
func (sr *ShapeReader) ReadByte() byte {
	var u uint8
	binary.Read(sr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt16 =
func (sr *ShapeReader) ReadUInt16() uint16 {
	var u uint16
	binary.Read(sr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt32 =
func (sr *ShapeReader) ReadUInt32() uint32 {
	var u uint32
	binary.Read(sr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt64 =
func (sr *ShapeReader) ReadUInt64() uint64 {
	var u uint64
	binary.Read(sr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUTF8 =
func (sr *ShapeReader) ReadUTF8(length uint32) string {
	buffer := sr.ReadBytes(length)
	return string(buffer)
}

// ReadGUID =
func (sr *ShapeReader) ReadGUID() Guid {
	buffer := sr.ReadBytes(16)
	return NewGuid(buffer)
}

// ReadString =
func (sr *ShapeReader) ReadString(maxLength uint32) string {
	buffer := make([]byte, 0, 128)
	i := uint32(0)
	for {
		c, _ := sr.reader.ReadByte()
		if c == 0 {
			break
		}
		buffer = append(buffer, c)
		i++
		if i+1 == maxLength {
			buffer[i+1] = 0
		}
	}

	return string(buffer)
}

// ReadCompressedUInt =
func (sr *ShapeReader) ReadCompressedUInt() uint32 {
	b1 := sr.ReadByte()
	if (b1 & 0x80) == 0x0 {
		return uint32(b1)
	}

	b2 := sr.ReadByte()
	if (b1 & 0xC0) == 0x80 {
		b1Shifted := uint16(b1&0x3f) << 8
		return uint32(b1Shifted & uint16(b2))
	}

	u1 := uint32(b1&0x3f) << 24
	u2 := uint32(b2) << 16
	u3 := uint32(sr.ReadByte()) << 8
	u4 := uint32(sr.ReadByte())
	return u1 | u2 | u3 | u4
}

// ReadBytes =
func (sr *ShapeReader) ReadBytes(length uint32) []byte {
	buffer := make([]byte, length)
	sr.reader.Read(buffer)
	return buffer
}

// Skip =
func (sr *ShapeReader) Skip(numBytes int32) {
	sr.reader.Seek(int64(numBytes), io.SeekCurrent)
}
