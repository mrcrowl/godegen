package reflect

import "godegen/cli"

type Field struct {
	name         string
	signature    *FieldSig
	memberAccess MemberAccess
	static       bool
	initOnly     bool
	literal      bool
}

const (
	FieldAttributesStatic   uint16 = 0x10
	FieldAttributesInitOnly uint16 = 0x20
	FieldAttributesLiteral  uint16 = 0x40
)

func newField(fieldRow *cli.FieldRow, asm *Assembly) *Field {
	sigBlob := fieldRow.GetSignatureBlob()
	sigReader := NewSignatureReader(sigBlob, asm)
	memberAccess := MemberAccess(fieldRow.Flags) & MemberAccessMask
	static := (fieldRow.Flags & FieldAttributesStatic) > 0
	initOnly := (fieldRow.Flags & FieldAttributesInitOnly) > 0
	literal := (fieldRow.Flags & FieldAttributesLiteral) > 0

	return &Field{
		name:         fieldRow.Name,
		signature:    sigReader.ReadFieldSignature(),
		memberAccess: memberAccess,
		static:       static,
		initOnly:     initOnly,
		literal:      literal,
	}
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Type() Type {
	return field.signature.fieldType
}
