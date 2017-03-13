package cli

type ModuleRefRow struct {
	Name string
	Mvid Guid
}

func (row *ModuleRefRow) String() string {
	return row.Name
}

func readModuleRefRow(sr *ShapeReader, streams *MetadataStreams, _ *TableSet) IRow {
	sr.Skip(2)
	name := streams.stringHeap.ReadString(sr)
	mvid := streams.guidHeap.ReadGuid(sr)
	sr.Skip(streams.stringHeap.GetIndexSizeInBytes())
	return &ModuleRefRow{
		Name: name,
		Mvid: mvid,
	}
}
