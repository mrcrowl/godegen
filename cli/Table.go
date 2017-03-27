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

func readInterfaceImpl(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	readSimpleIndex(sr, tables, TableIdxTypeDef)
	readTypeDefOrRefCodedIndex(sr, tables)
	return nil
}

func readMemberRef(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	readMemberRefParentCodedIndex(sr, tables)
	streams.stringHeap.ReadAndDiscard(sr)
	streams.blobHeap.ReadAndDiscard(sr)
	return nil
}

func readCustomAttribute(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	readHasCustomAttributeCodedIndex(sr, tables)
	readCustomAttributeTypeCodedIndex(sr, tables)
	streams.blobHeap.ReadAndDiscard(sr)
	return nil
}

func readFieldMarshal(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	readHasFieldMarshalCodedIndex(sr, tables)
	streams.blobHeap.ReadAndDiscard(sr)
	return nil
}

func readDeclSecurity(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	sr.ReadUInt16()
	readHasDeclSecurityCodedIndex(sr, tables)
	streams.blobHeap.ReadAndDiscard(sr)
	return nil
}

func readClassLayout(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	sr.ReadBytes(2)
	sr.ReadBytes(4)
	readSimpleIndex(sr, tables, TableIdxTypeDef)
	return nil
}

func readFieldLayout(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	sr.ReadBytes(4)
	readSimpleIndex(sr, tables, TableIdxField)
	return nil
}

func readStandAloneSig(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	streams.blobHeap.ReadAndDiscard(sr)
	return nil
}

func readEventMap(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	readSimpleIndex(sr, tables, TableIdxTypeDef)
	readSimpleIndex(sr, tables, TableIdxEvent)
	return nil
}

func readEvent(sr *ShapeReader, rowNumber uint32, streams *MetadataStreams, tables *TableSet) IRow {
	sr.ReadBytes(2)
	streams.stringHeap.ReadString(sr)
	readTypeDefOrRefCodedIndex(sr, tables)
	return nil
}

var rowReaderFns = [maxTableCount]RowReaderFn{
	0x00: readModuleRow,
	0x01: readTypeRefRow,
	0x02: readTypeDefRow,
	0x04: readFieldRow,
	0x06: readMethodDefRow,
	0x08: readParamRow,
	0x09: readInterfaceImpl,                            // InterfaceImpl
	0x0A: readMemberRef,                                // MemberRef
	0x0B: readConstRow,                                 // Const
	0x0C: readCustomAttribute,                          // CustomAttribute
	0x0D: readFieldMarshal,                             // FieldMarshal
	0x0E: readDeclSecurity,                             // DeclSecurity
	0x0F: readClassLayout,                              // ClassLayout
	0x10: readFieldLayout,                              // FieldLayout
	0x11: readStandAloneSig,                            // StandAloneSig
	0x12: readEventMap,                                 // EventMap
	0x14: readEvent,                                    // Event
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
