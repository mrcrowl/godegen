package cli

import "sort"

type Table struct {
	tableIndex uint8
	numRows    uint32
	sorted     bool
	rows       []IRow
}

type RowRange struct {
	from uint32
	to   uint32
}

func NewRowRange(from uint32, to uint32) RowRange {
	return RowRange{from, to}
}

func (rng RowRange) count() uint32 {
	return rng.to - rng.from
}

type RowReaderFn func(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow

func createPlaceholderReaderOfSize(nBytes uint32) RowReaderFn {
	return func(sr *ShapeReader, _ uint32, _ *MetadataStreams, _ *TableSet) IRow {
		sr.ReadBytes(nBytes)
		return nil
	}
}

var rowReaderFns = [maxTableCount]RowReaderFn{
	0x00: readModuleRow,
	0x01: readTypeRefRow,
	0x02: readTypeDefRow,
	0x04: readFieldRow,
	0x06: readMethodDefRow,
	0x08: readParamRow,
	0x09: createPlaceholderReaderOfSize(2 + 2),         // InterfaceImpl
	0x0A: createPlaceholderReaderOfSize(4 + 4 + 4),     // MemberRef
	0x0B: readConstRow,                                 // Const
	0x0C: createPlaceholderReaderOfSize(4 + 4 + 4),     // CustomAttribute
	0x0D: createPlaceholderReaderOfSize(4 + 4),         // FieldMarshal
	0x0E: createPlaceholderReaderOfSize(2 + 4 + 4),     // DeclSecurity
	0x0F: createPlaceholderReaderOfSize(2 + 4 + 2),     // ClassLayout
	0x10: createPlaceholderReaderOfSize(4 + 2),         // FieldLayout
	0x11: createPlaceholderReaderOfSize(4),             // StandAloneSig
	0x12: createPlaceholderReaderOfSize(2 + 2),         // EventMap
	0x14: createPlaceholderReaderOfSize(2 + 4 + 2),     // Event
	0x15: readPropertyMapRow,                           // PropertyMap
	0x17: readPropertyRow,                              // Property
	0x18: readMethodSemanticsRow,                       // MethodSemantics
	0x19: createPlaceholderReaderOfSize(2 + 2 + 2),     // MethodImpl
	0x1A: createPlaceholderReaderOfSize(4),             // ModuleRef
	0x1B: createPlaceholderReaderOfSize(4),             // TypeSpec
	0x1C: createPlaceholderReaderOfSize(2 + 4 + 4 + 2), // ImplMap
	0x1D: createPlaceholderReaderOfSize(4 + 2),         // FieldRVA
	0x20: readAssemblyRow,                              // Assembly
	0x21: createPlaceholderReaderOfSize(4),             // AssemblyProcessor
	0x22: createPlaceholderReaderOfSize(4 + 4 + 4),     // AssemblyOS
	0x23: readAssemblyRefRow,                           // AssemblyRef
	0x2A: readGenericParamRow,                          // GenericParam
	// createPlaceholderReaderOfSize(2 + 4 + 4),        // Constant
}

func newTable(tableIndex uint8, numRows uint32, isSorted bool) *Table {
	return &Table{
		tableIndex: tableIndex,
		numRows:    numRows,
		sorted:     isSorted,
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

func (table Table) ForEach(action func(IRow)) {
	for _, row := range table.rows {
		action(row)
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

func (table Table) BinarySearchRows(condition func(row IRow) bool) IRow {
	searchCriteriaFn := func(i int) bool {
		row := table.rows[i]
		return condition(row)
	}

	rowIndex := uint32(sort.Search(int(table.numRows), searchCriteriaFn))
	if rowIndex < table.numRows {
		return table.rows[rowIndex]
	}
	return nil
}
