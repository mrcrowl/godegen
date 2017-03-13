package cli

type TypeDefRow struct {
	Flags         uint32
	TypeName      string
	TypeNamespace string
	Extends       TypeDefOrRefIndex

	fieldRowRange  RowRange
	methodRowRange RowRange
}

type TypeDefOrRefType uint8

const (
	TDORTypeDef TypeDefOrRefType = iota
	TDORTypeRef
	TDORTypeSpec
)

type TypeDefOrRefIndex struct {
	Index uint32
	Type  TypeDefOrRefType
}

func NewTypeDefOrRefIndex(codedIndex CodedIndex) TypeDefOrRefIndex {
	return TypeDefOrRefIndex{
		Index: codedIndex.Index,
		Type:  TypeDefOrRefType(codedIndex.Tag),
	}
}

func (row *TypeDefRow) String() string {
	return row.TypeNamespace + "::" + row.TypeName
}

func readTypeDefRow(
	sr *ShapeReader,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	return &TypeDefRow{
		Flags:          sr.ReadUInt32(),
		TypeName:       streams.stringHeap.ReadString(sr),
		TypeNamespace:  streams.stringHeap.ReadString(sr),
		Extends:        NewTypeDefOrRefIndex(ReadCodedIndex(sr, tables, TableIdxTypeDef, TableIdxTypeRef, TableIdxTypeSpec)),
		fieldRowRange:  RowRange{from: ReadSimpleIndex(sr, tables, TableIdxField)},
		methodRowRange: RowRange{from: ReadSimpleIndex(sr, tables, TableIdxMethodDef)},
	}
}

func getTypeDefRow(rows []IRow, index uint32) *TypeDefRow {
	return rows[index].(*TypeDefRow)
}
