package cli

type ModuleRef struct {
	Name string
	Mvid Guid
}

func (row *ModuleRef) String() string {
	return row.Name
}

func readModuleRefRow(sr *ShapeReader, streams *MetadataStreams,
	_ *TableSet) IRow {
	sr.Skip(2)
	name := streams.stringHeap.ReadString(sr)
	mvid := streams.guidHeap.ReadGuid(sr)
	sr.Skip(streams.stringHeap.GetIndexSizeInBytes())
	return &ModuleRef{
		Name: name,
		Mvid: mvid,
	}
}
