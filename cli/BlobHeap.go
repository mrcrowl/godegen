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
	blob := heap.readBlobAtIndex(index)
	return blob
}

func (heap *BlobHeap) readBlobAtIndex(index uint32) *Blob {
	if err := heap.reader.Seek(int64(index)); err == nil {
		length := heap.reader.ReadCompressedUInt()
		return NewBlob(length, heap.reader.ReadBytes(length))
	}

	return ZeroBlob()
}
