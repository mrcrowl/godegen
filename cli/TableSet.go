package cli

import (
	"math"
)

const maxTableCount = 0x2d

type TableIdx uint8

const (
	TableIdxModule                 TableIdx = 0x00
	TableIdxTypeRef                TableIdx = 0x01
	TableIdxTypeDef                TableIdx = 0x02
	TableIdxField                  TableIdx = 0x04
	TableIdxMethodDef              TableIdx = 0x06
	TableIdxParam                  TableIdx = 0x08
	TableIdxInterfaceImpl          TableIdx = 0x09
	TableIdxMemberRef              TableIdx = 0x0A
	TableIdxConstant               TableIdx = 0x0B
	TableIdxCustomAttribute        TableIdx = 0x0C
	TableIdxFieldMarshal           TableIdx = 0x0D
	TableIdxDeclSecurity           TableIdx = 0x0E
	TableIdxClassLayout            TableIdx = 0x0F
	TableIdxFieldLayout            TableIdx = 0x10
	TableIdxStandAloneSig          TableIdx = 0x11
	TableIdxEventMap               TableIdx = 0x12
	TableIdxEvent                  TableIdx = 0x14
	TableIdxPropertyMap            TableIdx = 0x15
	TableIdxProperty               TableIdx = 0x17
	TableIdxMethodSemantics        TableIdx = 0x18
	TableIdxMethodImpl             TableIdx = 0x19
	TableIdxModuleRef              TableIdx = 0x1A
	TableIdxTypeSpec               TableIdx = 0x1B
	TableIdxImplMap                TableIdx = 0x1C
	TableIdxFieldRVA               TableIdx = 0x1D
	TableIdxAssembly               TableIdx = 0x20
	TableIdxAssemblyProcessor      TableIdx = 0x21
	TableIdxAssemblyOS             TableIdx = 0x22
	TableIdxAssemblyRef            TableIdx = 0x23
	TableIdxAssemblyRefProcessor   TableIdx = 0x24
	TableIdxAssemblyRefOS          TableIdx = 0x25
	TableIdxFile                   TableIdx = 0x26
	TableIdxExportedType           TableIdx = 0x27
	TableIdxManifestResource       TableIdx = 0x28
	TableIdxNestedClass            TableIdx = 0x29
	TableIdxGenericParam           TableIdx = 0x2A
	TableIdxMethodSpec             TableIdx = 0x2B
	TableIdxGenericParamConstraint TableIdx = 0x2C
)

type TableSet struct {
	tables []Table
}

func NewTableSet(streams *MetadataStreams) *TableSet {
	tableIndex := uint8(0)
	tildeStream := streams.tildeStream
	presentRowCounts, presentTables := tildeStream.RawRowCounts, tildeStream.ValidTables
	presentRowCountsIndex := 0
	tables := make([]Table, maxTableCount)
	for presentTables > 0 {
		isTablePresent := (presentTables & 0x1) > 0
		rowCount := uint32(0)
		if isTablePresent {
			rowCount = presentRowCounts[presentRowCountsIndex]
			presentRowCountsIndex++
		}

		table := newTable(tableIndex, rowCount)

		tables[tableIndex] = table
		tableIndex++
		presentTables >>= 1
	}

	return &TableSet{
		tables: tables,
	}
}

func (set *TableSet) GetRows(tableIndex TableIdx) []IRow {
	return set.tables[tableIndex].rows
}

func (set *TableSet) readAll(streams *MetadataStreams) {
	tablesReader := streams.tildeStream.GetTablesReader()
	for _, table := range set.tables {
		table.readRows(tablesReader, streams, set)
	}

	set.postRead()
}

func (set *TableSet) postRead() {
	set.collectTypeDefRanges()
	set.collectMethodDefParamRange()
}

func (set *TableSet) GetTable(index TableIdx) Table {
	return set.tables[index]
}

func (set *TableSet) collectTypeDefRanges() {
	typeDefTable := set.GetTable(TableIdxTypeDef)
	methodDefTable := set.GetTable(TableIdxMethodDef)
	fieldTable := set.GetTable(TableIdxField)
	rows := typeDefTable.rows
	numTypes := typeDefTable.numRows
	for i := uint32(0); i < numTypes; i++ {
		row := getTypeDefRow(rows, i)
		isLastRow := (i + 1) == numTypes
		if !isLastRow {
			nextRow := getTypeDefRow(rows, i+1)
			row.methodRowRange.to = nextRow.methodRowRange.from
			row.fieldRowRange.to = nextRow.fieldRowRange.from
		} else {
			row.methodRowRange.to = methodDefTable.numRows + 1
			row.fieldRowRange.to = fieldTable.numRows + 1
		}
	}
}

func (set *TableSet) collectMethodDefParamRange() {
	methodDefTable := set.GetTable(TableIdxMethodDef)
	paramTable := set.GetTable(TableIdxParam)
	rows := methodDefTable.rows
	numMethodRows := methodDefTable.numRows
	for i := uint32(0); i < numMethodRows; i++ {
		row := getMethodDefRow(rows, i)
		isLastRow := (i + 1) == numMethodRows
		if !isLastRow {
			nextRow := getMethodDefRow(rows, i+1)
			row.paramRowRange.to = nextRow.paramRowRange.from
		} else {
			row.paramRowRange.to = paramTable.numRows + 1
		}
	}
}

type tableReadInfo struct {
	UseBigIndex bool
	NumTagBits  uint8
}

func (set *TableSet) getTableReadInfo(tableIndexes []TableIdx) tableReadInfo {
	maxRowCount := uint32(0)
	numTables := float64(len(tableIndexes))
	bitsRequiredForTables := uint8(math.Ceil(math.Log2(numTables)))
	upperRowLimit := uint32(1 << (16 - bitsRequiredForTables))
	for _, tableIndex := range tableIndexes {
		rowCount := set.GetRowCount(tableIndex)
		if rowCount > maxRowCount {
			maxRowCount = rowCount
			if maxRowCount >= upperRowLimit {
				return tableReadInfo{NumTagBits: bitsRequiredForTables, UseBigIndex: true}
			}
		}
	}

	return tableReadInfo{
		NumTagBits:  bitsRequiredForTables,
		UseBigIndex: false,
	}
}

func (set *TableSet) useBigIndex(tableIndex TableIdx) bool {
	return set.tables[tableIndex].numRows >= 0x10000
}

func (set *TableSet) GetRowCount(tableIndex TableIdx) uint32 {
	return set.tables[tableIndex].numRows
}
