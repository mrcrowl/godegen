package reflect

import "educationperfect.com/godegen/cli"

type Field struct {
	name         string
	signature    *FieldSig
	memberAccess MemberAccess
	static       bool
	initOnly     bool
	literal      bool
	value        interface{}
}

const (
	FieldAttributesStatic     uint16 = 0x10
	FieldAttributesInitOnly   uint16 = 0x20
	FieldAttributesLiteral    uint16 = 0x40
	FieldAttributesHasDefault uint16 = 0x8000
)

func newField(fieldRow *cli.FieldRow, asm *Assembly) *Field {
	sigBlob := fieldRow.GetSignatureBlob()
	sigReader := NewSignatureReader(sigBlob, asm)
	memberAccess := MemberAccess(fieldRow.Flags) & MemberAccessMask
	static := (fieldRow.Flags & FieldAttributesStatic) > 0
	initOnly := (fieldRow.Flags & FieldAttributesInitOnly) > 0
	literal := (fieldRow.Flags & FieldAttributesLiteral) > 0
	hasDefault := (fieldRow.Flags & FieldAttributesHasDefault) > 0
	var value interface{}

	if hasDefault {
		constantTable := asm.metadata.Tables.GetTable(cli.TableIdxConstant)
		fieldRowNumber := fieldRow.RowNumber()
		selectedConstantIRow := constantTable.BinarySearchRows(func(row cli.IRow) bool {
			selectedConstantRow := row.(*cli.ConstRow)
			return selectedConstantRow.Parent.Row >= fieldRowNumber &&
				selectedConstantRow.Parent.Type >= cli.HCField
		})
		if selectedConstantIRow != nil {
			selectedConstantRow := selectedConstantIRow.(*cli.ConstRow)
			targetIndex := cli.HasConstantIndex{Row: fieldRowNumber, Type: cli.HCField}
			if selectedConstantRow.Parent == targetIndex {
				value = selectedConstantRow.ReadValue()
			}
		}
	}

	return &Field{
		name:         fieldRow.Name,
		signature:    sigReader.ReadFieldSignature(),
		memberAccess: memberAccess,
		static:       static,
		initOnly:     initOnly,
		literal:      literal,
		value:        value,
	}
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Type() Type {
	return field.signature.fieldType
}

func (field *Field) Value() interface{} {
	return field.value
}
