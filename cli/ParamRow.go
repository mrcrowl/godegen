package cli

type ParamRow struct {
	rowNumber uint32
	Flags     uint16
	sequence  uint16
	Name      string
}

func (row *ParamRow) String() string {
	return row.Name
}

func (row *ParamRow) RowNumber() uint32 {
	return row.rowNumber
}

func readParamRow(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	flags := sr.ReadUInt16()
	sequence := sr.ReadUInt16()
	name := streams.stringHeap.ReadString(sr)
	return &ParamRow{rowNumber, flags, sequence, name}
}

func getParamsInRange(rows []IRow, fromIndex uint32, toIndex uint32) []*ParamRow {
	numParams := toIndex - fromIndex
	params := make([]*ParamRow, numParams)
	selectedRows := rows[fromIndex:toIndex]
	for i, row := range selectedRows {
		params[i] = row.(*ParamRow)
	}
	return params
}
