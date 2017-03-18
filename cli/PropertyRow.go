package cli

type PropertyRow struct {
	rowNumber     uint32
	Flags         uint16
	Name          string
	signatureBlob Blob
}

func (row *PropertyRow) String() string {
	return row.Name
}

func (row *PropertyRow) RowNumber() uint32 {
	return row.rowNumber
}

func (row *PropertyRow) GetSignatureBlob() Blob {
	return row.signatureBlob
}

func readPropertyRow(
	sr *ShapeReader,
	rowNumber uint32,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	flags := sr.ReadUInt16()
	name := streams.stringHeap.ReadString(sr)
	signatureBlob := *streams.blobHeap.ReadBlob(sr)
	return &PropertyRow{rowNumber, flags, name, signatureBlob}
}

func getPropertyRowsInRange(rows []IRow, fromIndex uint32, toIndex uint32) []*PropertyRow {
	numProps := toIndex - fromIndex
	propRows := make([]*PropertyRow, numProps)
	selectedRows := rows[fromIndex:toIndex]
	for i, row := range selectedRows {
		propRows[i] = row.(*PropertyRow)
	}
	return propRows
}
