package cli

type PropertyMapRow struct {
	rowNumber        uint32
	Parent           uint32
	propertyRowRange RowRange
}

func (row *PropertyMapRow) RowNumber() uint32 {
	return row.rowNumber
}

func (row *PropertyMapRow) String() string {
	return ""
}

func readPropertyMapRow(
	sr *ShapeReader,
	rowNumber uint32,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	parentRowNumber := readSimpleIndex(sr, tables, TableIdxTypeDef)
	propertyStartRow := readSimpleIndex(sr, tables, TableIdxProperty)
	return &PropertyMapRow{rowNumber, parentRowNumber, RowRange{propertyStartRow, propertyStartRow}}
}

func getPropertyMapRow(rows []IRow, index uint32) *PropertyMapRow {
	return rows[index].(*PropertyMapRow)
}
