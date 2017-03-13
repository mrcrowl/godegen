package cli

type FieldRow struct {
	Flags         uint16
	Name          string
	signatureBlob Blob
}

func (row *FieldRow) String() string {
	return row.Name
}

func readFieldRow(
	sr *ShapeReader,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	flags := sr.ReadUInt16()
	name := streams.stringHeap.ReadString(sr)
	signatureBlob := *streams.blobHeap.ReadBlob(sr)
	return &FieldRow{flags, name, signatureBlob}
}
