package cli

import (
	"bytes"
)

type MethodDefRow struct {
	RVA           uint32
	ImplFlags     uint16
	Flags         uint16
	Name          string
	signatureBlob Blob
	paramRowRange RowRange
	paramTable    *Table
}

func (row *MethodDefRow) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(row.Name)
	buf.WriteByte('(')
	numRows := row.paramRowRange.count()
	if numRows > 0 {
		startIndex := row.paramRowRange.from - 1
		endIndex := row.paramRowRange.to - 1
		params := getParamsInRange(row.paramTable.rows, startIndex, endIndex)
		for i, param := range params {
			buf.WriteString(param.Name)
			if uint32(i+1) < numRows {
				buf.WriteString(", ")
			}
		}
	}
	buf.WriteByte(')')
	return buf.String()
}

func readMethodDefRow(sr *ShapeReader, streams *MetadataStreams, tables *TableSet) IRow {
	rva := sr.ReadUInt32()
	implFlags := sr.ReadUInt16()
	flags := sr.ReadUInt16()
	name := streams.stringHeap.ReadString(sr)
	signatureBlob := *streams.blobHeap.ReadBlob(sr)
	paramFromIndex := ReadSimpleIndex(sr, tables, TableIdxMethodDef)
	return &MethodDefRow{rva, implFlags, flags, name, signatureBlob, RowRange{paramFromIndex, paramFromIndex}, nil}
}

func getMethodDefRow(rows []IRow, index uint32) *MethodDefRow {
	return rows[index].(*MethodDefRow)
}
