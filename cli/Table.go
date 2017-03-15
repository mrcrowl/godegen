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

type RowReaderFn func(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow

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
		table.rows[i] = readerFn(tr, i+1, streams, tables)
	}
}

func (table Table) Where(condition func(IRow) bool) []IRow {
	matches := make([]IRow, 0, 128)
	for _, row := range table.rows {
		if condition(row) {
			matches = append(matches, row)
		}
	}
	return matches
}

func (table Table) First(condition func(IRow) bool) IRow {
	for _, row := range table.rows {
		if condition(row) {
			return row
		}
	}
	return nil
}

func (table Table) RowNumberWhere(condition func(IRow) bool) uint32 {
	for i, row := range table.rows {
		if condition(row) {
			return uint32(i + 1)
		}
	}
	return 0
}

func (table Table) GetRow(rowNumber uint32) IRow {
	if rowNumber > 0 && rowNumber <= table.numRows {
		return table.rows[rowNumber-1]
	}
	return nil
}
