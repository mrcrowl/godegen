package cli

type Version struct {
	major    uint16
	minor    uint16
	build    uint16
	revision uint16
}

type AssemblyRefRow struct {
	rowNumber        uint32
	version          Version
	flags            uint32
	publicKeyOrToken Blob
	Name             string
	Culture          string
	HashValue        Blob
}

func (row *AssemblyRefRow) String() string {
	return row.Name
}

func (row *AssemblyRefRow) RowNumber() uint32 {
	return row.rowNumber
}

func readAssemblyRefRow(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, _ *TableSet) IRow {
	return &AssemblyRefRow{
		rowNumber: rowNumber,
		version: Version{
			sr.ReadUInt16(),
			sr.ReadUInt16(),
			sr.ReadUInt16(),
			sr.ReadUInt16(),
		},
		flags:            sr.ReadUInt32(),
		publicKeyOrToken: *streams.blobHeap.ReadBlob(sr),
		Name:             streams.stringHeap.ReadString(sr),
		Culture:          streams.stringHeap.ReadString(sr),
		HashValue:        *streams.blobHeap.ReadBlob(sr),
	}
}
