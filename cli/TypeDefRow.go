package cli

type TypeDefRow struct {
	rowNumber     uint32
	Flags         uint32
	TypeName      string
	TypeNamespace string
	ExtendsIndex  TypeDefOrRefIndex

	fieldRowRange  RowRange
	methodRowRange RowRange
}

func (row *TypeDefRow) RowNumber() uint32 {
	return row.rowNumber
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
	methods := getMethodsInRange(rows, startIndex, endIndex)
	return methods
}

func (row *TypeDefRow) GetFieldRows(set *TableSet) []*FieldRow {
	rowRange := row.fieldRowRange
	startIndex := rowRange.from - 1
	endIndex := rowRange.to - 1
	rows := set.GetTable(TableIdxField).rows
	fields := getFieldsInRange(rows, startIndex, endIndex)
	return fields
}

func (row *TypeDefRow) GetPropertyRows(set *TableSet) []*PropertyRow {
	// find the property map row
	propertyMapTable := set.GetTable(TableIdxPropertyMap)
	typeDefRowNumber := row.RowNumber()
	selectedMapIRow := propertyMapTable.BinarySearchRows(func(row IRow) bool {
		propertyMapRow := row.(*PropertyMapRow)
		return propertyMapRow.Parent >= typeDefRowNumber
	})
	if selectedMapIRow == nil || selectedMapIRow.(*PropertyMapRow).Parent != typeDefRowNumber {
		return []*PropertyRow{}
	}

	// select the property rows (using the property map's range)
	selectedMapRow := selectedMapIRow.(*PropertyMapRow)
	rowRange := selectedMapRow.propertyRowRange
	startIndex := rowRange.from - 1
	endIndex := rowRange.to - 1
	propertyRows := set.GetTable(TableIdxProperty).rows
	properties := getPropertyRowsInRange(propertyRows, startIndex, endIndex)

	return properties
}

func readTypeDefRow(
	sr *ShapeReader,
	rowNumber uint32,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	flags := sr.ReadUInt32()
	typeName := streams.stringHeap.ReadString(sr)
	typeNamespace := streams.stringHeap.ReadString(sr)
	codedIndex := readCodedIndex(sr, tables, TableIdxTypeDef, TableIdxTypeRef, TableIdxTypeSpec)
	fieldFrom := readSimpleIndex(sr, tables, TableIdxField)
	methodFrom := readSimpleIndex(sr, tables, TableIdxMethodDef)
	return &TypeDefRow{
		rowNumber:      rowNumber,
		Flags:          flags,
		TypeName:       typeName,
		TypeNamespace:  typeNamespace,
		ExtendsIndex:   newTypeDefOrRefIndex(codedIndex),
		fieldRowRange:  RowRange{from: fieldFrom},
		methodRowRange: RowRange{from: methodFrom},
	}
}

func getTypeDefRow(rows []IRow, index uint32) *TypeDefRow {
	return rows[index].(*TypeDefRow)
}
