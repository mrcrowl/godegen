package cli

type ModuleRefRow struct {
	rowNumber uint32
	Name      string
	Mvid      Guid
}

func (row *ModuleRefRow) String() string {
	return row.Name
}

func (row *ModuleRefRow) RowNumber() uint32 {
	return row.rowNumber
}

func readModuleRefRow(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, _ *TableSet) IRow {
	sr.Skip(2)
	name := streams.stringHeap.ReadString(sr)
	mvid := streams.guidHeap.ReadGuid(sr)
	sr.Skip(streams.stringHeap.GetIndexSizeInBytes())
	return &ModuleRefRow{
		Name:      name,
		Mvid:      mvid,
		rowNumber: rowNumber,
	}
}
