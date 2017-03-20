package cli

import (
	"bytes"
)

type Blob struct {
	Length uint32
	Data   []byte
}

const (
	ELEMENT_TYPE_VOID    = 0x01
	ELEMENT_TYPE_BOOLEAN = 0x02
	ELEMENT_TYPE_CHAR    = 0x03
	ELEMENT_TYPE_I1      = 0x04
	ELEMENT_TYPE_U1      = 0x05
	ELEMENT_TYPE_I2      = 0x06
	ELEMENT_TYPE_U2      = 0x07
	ELEMENT_TYPE_I4      = 0x08
	ELEMENT_TYPE_U4      = 0x09
	ELEMENT_TYPE_I8      = 0x0a
	ELEMENT_TYPE_U8      = 0x0b
	ELEMENT_TYPE_R4      = 0x0c
	ELEMENT_TYPE_R8      = 0x0d
	ELEMENT_TYPE_STRING  = 0x0e
)

func NewBlob(length uint32, buffer []byte) *Blob {
	return &Blob{length, buffer}
}

func ZeroBlob() *Blob {
	return &Blob{}
}

func (blob *Blob) ReadTypedData(typeID byte) interface{} {
	sr := NewShapeReader(bytes.NewReader(blob.Data))
	switch typeID {
	case ELEMENT_TYPE_VOID:
		return nil
	case ELEMENT_TYPE_BOOLEAN:
		return sr.ReadBoolean()
	case ELEMENT_TYPE_CHAR:
		return sr.ReadUTF16(2)
	case ELEMENT_TYPE_I1:
		return sr.ReadByte()
	case ELEMENT_TYPE_U1:
		return sr.ReadInt8()
	case ELEMENT_TYPE_I2:
		return sr.ReadInt16()
	case ELEMENT_TYPE_U2:
		return sr.ReadUInt16()
	case ELEMENT_TYPE_I4:
		return sr.ReadInt32()
	case ELEMENT_TYPE_U4:
		return sr.ReadUInt64()
	case ELEMENT_TYPE_I8:
		return sr.ReadInt64()
	case ELEMENT_TYPE_U8:
		return sr.ReadUInt64()
	case ELEMENT_TYPE_R4:
		return sr.ReadFloat32()
	case ELEMENT_TYPE_R8:
		return sr.ReadFloat64()
	case ELEMENT_TYPE_STRING:
		return sr.ReadUTF16(blob.Length)
	}

	return nil
}
