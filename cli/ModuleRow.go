package cli

type ModuleRow struct {
	rowNumber uint32
	Name      string
	Mvid      Guid
}

func (row *ModuleRow) String() string {
	return row.Name
}

func (row *ModuleRow) RowNumber() uint32 {
	return row.rowNumber
}

func readModuleRow(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, _ *TableSet) IRow {
	sr.Skip(2)
	name := streams.stringHeap.ReadString(sr)
	mvid := streams.guidHeap.ReadGuid(sr)
	sr.Skip(streams.guidHeap.GetIndexSizeInBytes() * 2)
	return &ModuleRow{
		Name:      name,
		Mvid:      mvid,
		rowNumber: rowNumber,
	}
}
