package cli

type Table struct {
	tableIndex uint8
	numRows    uint32
	rows       []IRow
}

type RowRange struct {
	from uint32
	to   uint32
}

func (rng RowRange) count() uint32 {
	return rng.to - rng.from
}

type RowReaderFn func(sr *ShapeReader, streams *MetadataStreams, tables *TableSet) IRow

var rowReaderFns = [maxTableCount]RowReaderFn{
	readModuleRefRow,
	readTypeRefRow,
	readTypeDefRow,
	nil,
	readFieldRow,
	nil,
	readMethodDefRow,
	nil,
	readParamRow,
}

func newTable(tableIndex uint8, numRows uint32) Table {
	return Table{
		tableIndex: tableIndex,
		numRows:    numRows,
		rows:       make([]IRow, numRows),
	}
}

func (table Table) readRows(tr *ShapeReader, streams *MetadataStreams, tables *TableSet) {
	readerFn := rowReaderFns[table.tableIndex]
	if readerFn == nil {
		return
	}

	for i := uint32(0); i < table.numRows; i++ {
		table.rows[i] = readerFn(tr, streams, tables)
	}
}
