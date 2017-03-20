package cli

type ConstRow struct {
	rowNumber uint32
	TypeID    byte
	Parent    HasConstantIndex
	ValueBlob Blob
}

func (row *ConstRow) RowNumber() uint32 {
	return row.rowNumber
}

func (row *ConstRow) String() string {
	return ""
}

func (row *ConstRow) ReadValue() interface{} {
	value := row.ValueBlob.ReadTypedData(row.TypeID)
	return value
}

func readConstRow(
	sr *ShapeReader,
	rowNumber uint32,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	typeID := sr.ReadByte()
	sr.ReadByte() // padding
	parent := newHasConstantIndex(readCodedIndex(sr, tables, TableIdxField, TableIdxParam, TableIdxProperty))
	value := *streams.blobHeap.ReadBlob(sr)

	return &ConstRow{rowNumber, typeID, parent, value}
}

// func getPropertyMapRow(rows []IRow, index uint32) *PropertyMapRow {
// 	return rows[index].(*PropertyMapRow)
// }
