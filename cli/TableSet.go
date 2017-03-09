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

func (set *TableSet) GetRows(tableIndex TableIdx) []IRow {
	return set.tables[tableIndex].rows
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

		table := NewTable(tableIndex, rowCount)

		tables[tableIndex] = table
		tableIndex++
		presentTables >>= 1
	}

	return &TableSet{
		tables: tables,
	}
}

func (set *TableSet) ReadAll(streams *MetadataStreams) {
	tablesReader := streams.tildeStream.GetTablesReader()
	for _, table := range set.tables {
		table.ReadRows(tablesReader, streams, set)
	}
}

type TableReadInfo struct {
	UseBigIndex bool
	NumTagBits  uint8
}

func (tables *TableSet) GetTableReadInfo(tableIndexes []TableIdx) TableReadInfo {
	maxRowCount := uint32(0)
	numTables := float64(len(tableIndexes))
	bitsRequiredForTables := uint8(math.Ceil(math.Log2(numTables)))
	upperRowLimit := uint32(1 << (16 - bitsRequiredForTables))
	for _, tableIndex := range tableIndexes {
		rowCount := tables.GetRowCount(tableIndex)
		if rowCount > maxRowCount {
			maxRowCount = rowCount
			if maxRowCount >= upperRowLimit {
				return TableReadInfo{NumTagBits: bitsRequiredForTables, UseBigIndex: true}
			}
		}
	}

	return TableReadInfo{
		NumTagBits:  bitsRequiredForTables,
		UseBigIndex: false,
	}
}

func (tables *TableSet) GetRowCount(tableIndex TableIdx) uint32 {
	return tables.tables[tableIndex].numRows
}
