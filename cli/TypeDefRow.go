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

// public static async Task<TasksForUserReturnObject> GetCurrentTasksForUser(LPSession session)
// {13 [0 1 21 18 128 225 1 18 148 164 18 145 188]}
//      DEFAULT   TR: TASK  CLASS      CLASS
//        ONE PARAM       GENARGS		  TD: 3876---
//          GENINST            TD: 529
//             CLASS

func NewTypeDefOrRefIndex(codedIndex CodedIndex) TypeDefOrRefIndex {
	return TypeDefOrRefIndex{
		Index: codedIndex.Index,
		Type:  TypeDefOrRefType(codedIndex.Tag),
	}
}

func (row *TypeDefRow) String() string {
	return row.FullName()
}

func (row *TypeDefRow) FullName() string {
	return row.TypeNamespace + "." + row.TypeName
}

func (row *TypeDefRow) GetMethodRows(set *TableSet) []*MethodDefRow {
	rowRange := row.methodRowRange
	startIndex := rowRange.from - 1
	endIndex := rowRange.to - 1
	rows := set.GetTable(TableIdxMethodDef).rows
	params := getMethodsInRange(rows, startIndex, endIndex)
	return params
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
