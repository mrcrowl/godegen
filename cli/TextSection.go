package cli

import (
	"bytes"
	"debug/pe"
)

// TextSection wraps the .text section of a PE file
type TextSection struct {
	buffer []byte
	rva    uint32
	size   uint32
}

// NewTextSection is
func NewTextSection(section *pe.Section) *TextSection {
	textSection := &TextSection{
		buffer: make([]byte, section.Size),
		rva:    section.VirtualAddress,
		size:   section.Size,
	}
	reader := section.Open()
	reader.Read(textSection.buffer)
	return textSection
}

// GetRange =
func (ts *TextSection) GetRange(rva uint32, length uint32) *bytes.Buffer {
	offset := rva - ts.rva
	return bytes.NewBuffer(ts.buffer[offset : offset+length])
}

// GetReaderAt =
func (ts *TextSection) GetReaderAt(rva uint32) *bytes.Reader {
	offset := rva - ts.rva
	return bytes.NewReader(ts.buffer[offset:])
}
