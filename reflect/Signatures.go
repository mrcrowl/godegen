package reflect

import "godegen/cli"
import "bytes"

const (
	DEFAULT      = 0x0
	VARARG       = 0x5
	GENERIC      = 0x10
	HASTHIS      = 0x20
	EXPLICITTHIS = 0x40
)

const (
	ELEMENT_TYPE_END         = 0x00
	ELEMENT_TYPE_VOID        = 0x01
	ELEMENT_TYPE_BOOLEAN     = 0x02
	ELEMENT_TYPE_CHAR        = 0x03
	ELEMENT_TYPE_I1          = 0x04
	ELEMENT_TYPE_U1          = 0x05
	ELEMENT_TYPE_I2          = 0x06
	ELEMENT_TYPE_U2          = 0x07
	ELEMENT_TYPE_I4          = 0x08
	ELEMENT_TYPE_U4          = 0x09
	ELEMENT_TYPE_I8          = 0x0a
	ELEMENT_TYPE_U8          = 0x0b
	ELEMENT_TYPE_R4          = 0x0c
	ELEMENT_TYPE_R8          = 0x0d
	ELEMENT_TYPE_STRING      = 0x0e
	ELEMENT_TYPE_PTR         = 0x0f
	ELEMENT_TYPE_BYREF       = 0x10
	ELEMENT_TYPE_VALUETYPE   = 0x11
	ELEMENT_TYPE_CLASS       = 0x12
	ELEMENT_TYPE_VAR         = 0x13
	ELEMENT_TYPE_ARRAY       = 0x14
	ELEMENT_TYPE_GENERICINST = 0x15
	ELEMENT_TYPE_TYPEDBYREF  = 0x16
	ELEMENT_TYPE_I           = 0x18
	ELEMENT_TYPE_U           = 0x19
	ELEMENT_TYPE_FNPTR       = 0x1b
	ELEMENT_TYPE_OBJECT      = 0x1c
	ELEMENT_TYPE_SZARRAY     = 0x1d
	ELEMENT_TYPE_MVAR        = 0x1e
	ELEMENT_TYPE_CMOD_REQD   = 0x1f
	ELEMENT_TYPE_CMOD_OPT    = 0x20
	ELEMENT_TYPE_INTERNAL    = 0x21
	ELEMENT_TYPE_MODIFIER    = 0x40
	ELEMENT_TYPE_SENTINEL    = 0x41
	ELEMENT_TYPE_PINNED      = 0x45
	ELEMENT_TYPE_TYPE        = 0x50
)

type SignatureReader struct {
	assembly *Assembly
	shape    *cli.ShapeReader
	// rawBlob cli.Blob
}

func NewSignatureReader(blob cli.Blob, assembly *Assembly) *SignatureReader {
	sr := cli.NewShapeReader(bytes.NewReader(blob.Data))
	return &SignatureReader{assembly, sr} //, blob}
}

type MethodSig struct {
	returnType       Type
	parameters       []*Parameter
	hasThis          bool
	explicitThisType Type
}

func (sig *SignatureReader) ReadMethodSignature(paramRows []*cli.ParamRow) *MethodSig {
	flags := sig.shape.ReadByte()
	hasThis := (flags & HASTHIS) > 0
	hasExplicitThis := (flags & EXPLICITTHIS) > 0
	// vararg := (flags & VARARG) > 0
	generic := (flags & GENERIC) > 0
	// genParamCount := 0
	if generic {
		// genParamCount := sig.shape.ReadCompressedUInt()
		sig.shape.ReadCompressedUInt()
	}

	paramCount := sig.shape.ReadCompressedUInt()

	var explicitThisType Type
	if hasThis && hasExplicitThis {
		thisParam := sig.ReadParam("this")
		explicitThisType = thisParam.Type()
	}

	typeByte := sig.shape.ReadByte()
	returnType := sig.ReadTypeWithID(typeByte)

	params := make([]*Parameter, paramCount)
	for i := uint32(0); i < paramCount; i++ {
		name := paramRows[i].Name
		params[i] = sig.ReadParam(name)
	}

	return &MethodSig{
		returnType,
		params,
		hasThis,
		explicitThisType,
	}
}

type FieldSig struct {
	fieldType Type
}

func (sig *SignatureReader) ReadFieldSignature() *FieldSig {
	sig.shape.ReadByte()
	fieldType := sig.ReadType()

	return &FieldSig{
		fieldType,
	}
}

func (sig *SignatureReader) ReadType() Type {
	id := sig.shape.ReadByte()
	return sig.ReadTypeWithID(id)
}

func (sig *SignatureReader) ReadTypeWithID(id byte) Type {
	builtIn := sig.assembly.typeCache.getBuiltIn(id)
	if builtIn != nil {
		return builtIn
	}

	switch id {
	case ELEMENT_TYPE_CLASS:
		return sig.ReadClassType()

	case ELEMENT_TYPE_VALUETYPE:
		return sig.ReadValueType()

	case ELEMENT_TYPE_GENERICINST:
		return sig.ReadGenericInstType()
	}
	return nil
}

func (sig *SignatureReader) ReadClassType() Type {
	typeIndex := sig.ReadTypeDefOrRefOrSpecEncoded()
	typ := sig.assembly.getTypeByIndex(typeIndex)
	return typ
}

func (sig *SignatureReader) ReadValueType() Type {
	return sig.ReadClassType()
}

func (sig *SignatureReader) ReadGenericInstType() Type {
	sig.shape.ReadByte()
	typeIndex := sig.ReadTypeDefOrRefOrSpecEncoded()
	typ := sig.assembly.getTypeByIndex(typeIndex)
	genArgCount := sig.shape.ReadCompressedUInt()
	genArgs := make([]Type, genArgCount)
	for i := uint32(0); i < genArgCount; i++ {
		genArgs[i] = sig.ReadType()
	}
	return newGenericType(typ, genArgs, sig.assembly)
}

func (sig *SignatureReader) ReadParam(name string) *Parameter {
	var byref bool
	var typ Type

	b := sig.shape.ReadByte()
	for isCustomMod(b) {
		sig.ReadTypeDefOrRefOrSpecEncoded()
		b = sig.shape.ReadByte()
	}

	if b == ELEMENT_TYPE_TYPEDBYREF {
		typ = sig.assembly.typeCache.builtInTypes[ELEMENT_TYPE_TYPEDBYREF]
	} else {
		if b == ELEMENT_TYPE_BYREF {
			byref = true
			b = sig.shape.ReadByte()
		}

		typ = sig.ReadTypeWithID(b)
	}

	if byref {

	}

	return newParameter(name, typ)
}

func (sig *SignatureReader) ReadTypeDefOrRefOrSpecEncoded() cli.TypeDefOrRefIndex {
	codedIndex := sig.shape.ReadCompressedUInt()
	table := codedIndex & 0x3
	row := codedIndex >> 2
	typ := cli.TypeDefOrRefType(uint8(table))
	return cli.TypeDefOrRefIndex{row, typ}
}

func isCustomMod(b byte) bool {
	return b == ELEMENT_TYPE_CMOD_OPT || b == ELEMENT_TYPE_CMOD_REQD
}
