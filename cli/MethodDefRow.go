package cli

type MethodDefRow struct {
	RVA           uint32
	ImplFlags     uint16
	Flags         uint16
	Name          string
	signatureBlob Blob
	paramRowRange RowRange
}

func (row *MethodDefRow) String() string {
	return row.Name
}

func (row *MethodDefRow) GetSignature() Blob {
	return row.signatureBlob
}

func (row *MethodDefRow) GetParams(set *TableSet) []*ParamRow {
	rowRange := row.paramRowRange
	startIndex := rowRange.from - 1
	endIndex := rowRange.to - 1
	rows := set.GetTable(TableIdxParam).rows
	params := getParamsInRange(rows, startIndex, endIndex)
	return params
}

func readMethodDefRow(sr *ShapeReader, streams *MetadataStreams, tables *TableSet) IRow {
	rva := sr.ReadUInt32()
	implFlags := sr.ReadUInt16()
	flags := sr.ReadUInt16()
	name := streams.stringHeap.ReadString(sr)
	signatureBlob := *streams.blobHeap.ReadBlob(sr)
	paramFromIndex := ReadSimpleIndex(sr, tables, TableIdxMethodDef)
	paramRowRange := RowRange{paramFromIndex, paramFromIndex}
	return &MethodDefRow{rva, implFlags, flags, name, signatureBlob, paramRowRange}
}

func getMethodDefRow(rows []IRow, index uint32) *MethodDefRow {
	return rows[index].(*MethodDefRow)
}

func getMethodsInRange(rows []IRow, fromIndex uint32, toIndex uint32) []*MethodDefRow {
	numParams := toIndex - fromIndex
	methodRows := make([]*MethodDefRow, numParams)
	selectedRows := rows[fromIndex:toIndex]
	for i, row := range selectedRows {
		methodRows[i] = row.(*MethodDefRow)
	}
	return methodRows
}
