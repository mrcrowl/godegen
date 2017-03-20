package cli

type GenericParamRow struct {
	rowNumber uint32
	number    uint16
	flags     uint16
	owner     TypeOrMethodDefIndex
	Name      string
}

func (row *GenericParamRow) String() string {
	return row.Name
}

func (row *GenericParamRow) RowNumber() uint32 {
	return row.rowNumber
}

func readGenericParamRow(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	return &GenericParamRow{
		rowNumber: rowNumber,
		number:    sr.ReadUInt16(),
		flags:     sr.ReadUInt16(),
		owner:     newTypeOrMethodDefIndex(readCodedIndex(sr, tables, TableIdxTypeDef, TableIdxMethodDef)),
		Name:      streams.stringHeap.ReadString(sr),
	}
}
