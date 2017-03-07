package cli

import (
	"bytes"
	"encoding/binary"
	"io"
)

type MetaReader struct {
	reader *bytes.Reader
}

func NewMetaReader(reader *bytes.Reader) *MetaReader {
	metaReader := &MetaReader{
		reader: reader,
	}
	return metaReader
}

// ReadUInt16 =
func (mr *MetaReader) ReadUInt16() uint16 {
	var u uint16
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt32 =
func (mr *MetaReader) ReadUInt32() uint32 {
	var u uint32
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUInt64 =
func (mr *MetaReader) ReadUInt64() uint64 {
	var u uint64
	binary.Read(mr.reader, binary.LittleEndian, &u)
	return u
}

// ReadUTF8 =
func (mr *MetaReader) ReadUTF8(length uint32) string {
	buffer := make([]byte, length)
	mr.reader.Read(buffer)
	return string(buffer)
}

// Skip =
func (mr *MetaReader) Skip(numBytes int32) {
	mr.reader.Seek(int64(numBytes), io.SeekCurrent)
}
