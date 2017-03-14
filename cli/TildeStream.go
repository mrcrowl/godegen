package cli

type HeapFlag byte

const (
	StringHeapFlag HeapFlag = 0x1
	GuidHeapFlag   HeapFlag = 0x2
	BlobHeapFlag   HeapFlag = 0x3
)

type TildeStream struct {
	heapSizes    HeapFlag
	ValidTables  uint64
	sortedTables uint64
	numTables    byte
	RawRowCounts []uint32
	tablesReader *ShapeReader
}

// NewTildeStream =
func NewTildeStream(sr *ShapeReader) *TildeStream {
	tilde := new(TildeStream)
	sr.Skip(6)
	tilde.heapSizes = HeapFlag(sr.ReadByte())
	sr.Skip(1)
	tilde.ValidTables = sr.ReadUInt64()
	tilde.sortedTables = sr.ReadUInt64()
	tilde.numTables = countSetBits(tilde.ValidTables)
	rowCounts := make([]uint32, tilde.numTables)
	for i := uint8(0); i < tilde.numTables; i++ {
		rowCounts[i] = sr.ReadUInt32()
	}
	tilde.RawRowCounts = rowCounts
	tilde.tablesReader = sr
	return tilde
}

func (tilde *TildeStream) GetTablesReader() *ShapeReader {
	return tilde.tablesReader
}

func (tilde *TildeStream) doesHeapUseBigIndex(heap HeapFlag) bool {
	return (tilde.heapSizes & heap) > 0
}

func countSetBits(vector uint64) uint8 {
	bits := uint8(0)
	for vector > 0 {
		if (vector & 0x1) > 0 {
			bits++
		}
		vector >>= 1
	}
	return bits
}
