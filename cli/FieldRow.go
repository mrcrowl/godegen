package cli

type FieldRow struct {
	rowNumber     uint32
	Flags         uint16
	Name          string
	signatureBlob Blob
}

func (row *FieldRow) String() string {
	return row.Name
}

func (row *FieldRow) RowNumber() uint32 {
	return row.rowNumber
}

func (row *FieldRow) GetSignatureBlob() Blob {
	return row.signatureBlob
}

func readFieldRow(
	sr *ShapeReader,
	rowNumber uint32,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	flags := sr.ReadUInt16()
	name := streams.stringHeap.ReadString(sr)
	signatureBlob := *streams.blobHeap.ReadBlob(sr)
	return &FieldRow{rowNumber, flags, name, signatureBlob}
}
