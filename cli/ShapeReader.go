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

// ReadUInt8 =
func (mr *ShapeReader) ReadUInt8() uint8 {
	var u uint8
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt16 =
func (mr *ShapeReader) ReadUInt16() uint16 {
	var u uint16
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt32 =
func (mr *ShapeReader) ReadUInt32() uint32 {
	var u uint32
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt64 =
func (mr *ShapeReader) ReadUInt64() uint64 {
	var u uint64
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUTF8 =
func (mr *ShapeReader) ReadUTF8(length uint32) string {
	buffer := mr.ReadBytes(length)
	return string(buffer)
}

// ReadGuid =
func (mr *ShapeReader) ReadGuid() Guid {
	buffer := mr.ReadBytes(16)
	return NewGuid(buffer)
}

// ReadString =
func (mr *ShapeReader) ReadString(maxLength uint32) string {
	buffer := make([]byte, 0, 128)
	i := uint32(0)
	for {
		c, _ := mr.reader.ReadByte()
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

// ReadBytes =
func (mr *ShapeReader) ReadBytes(length uint32) []byte {
	buffer := make([]byte, length)
	mr.reader.Read(buffer)
	return buffer
}

// Skip =
func (mr *ShapeReader) Skip(numBytes int32) {
	mr.reader.Seek(int64(numBytes), io.SeekCurrent)
}
