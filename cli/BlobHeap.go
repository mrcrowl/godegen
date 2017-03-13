package cli

type BlobHeap struct {
	Heap
}

// NewBlobHeap =
func NewBlobHeap(sr *ShapeReader, useBigIndex bool) *BlobHeap {
	heap := &BlobHeap{
		Heap{
			useBigIndex: useBigIndex,
			reader:      sr,
		},
	}

	return heap
}

func (heap *BlobHeap) ReadBlob(sr *ShapeReader) *Blob {
	var index uint32
	if heap.useBigIndex {
		index = sr.ReadUInt32()
	} else {
		index = uint32(sr.ReadUInt16())
	}
	blob := heap.ReadBlobAtIndex(index)
	return blob
}

func (heap *BlobHeap) ReadBlobAtIndex(index uint32) *Blob {
	if err := heap.reader.Seek(int64(index)); err == nil {
		var length uint32 = readBlobLength(heap.reader)
		return NewBlob(length, heap.reader.ReadBytes(length))
	}

	return ZeroBlob()
}

func readBlobLength(sr *ShapeReader) uint32 {
	// single byte?
	byte1 := sr.ReadUInt8()
	if (byte1 & 0x80) == 0 {
		return uint32(byte1)
	}

	// two bytes?
	byte2 := sr.ReadUInt8()
	if (byte1 & 0xC0) == 0x80 {
		return uint32((byte1&0x3F)<<8 + byte2)
	}

	// four bytes?
	byte3 := sr.ReadUInt8()
	byte4 := sr.ReadUInt8()
	if (byte1 & 0xE0) == 0xC0 {
		return uint32((byte1&0x1F)<<24 + byte2<<16 + byte3<<8 + byte4)
	}

	panic("Invalid blob length")
}
