package cli

type AssemblyRow struct {
	rowNumber uint32
	version   Version
	flags     uint32
	publicKey Blob
	Name      string
	Culture   string
}

func (row *AssemblyRow) String() string {
	return row.Name
}

func (row *AssemblyRow) RowNumber() uint32 {
	return row.rowNumber
}

func readAssemblyRow(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, _ *TableSet) IRow {
	sr.Skip(4)
	return &AssemblyRow{
		rowNumber: rowNumber,
		version: Version{
			sr.ReadUInt16(),
			sr.ReadUInt16(),
			sr.ReadUInt16(),
			sr.ReadUInt16(),
		},
		flags:     sr.ReadUInt32(),
		publicKey: *streams.blobHeap.ReadBlob(sr),
		Name:      streams.stringHeap.ReadString(sr),
		Culture:   streams.stringHeap.ReadString(sr),
	}
}
