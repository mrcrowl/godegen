package cli

type GuidHeap struct {
	Heap
}

// NewTildeStream =
func NewGuidHeap(sr *ShapeReader, useBigIndex bool) *GuidHeap {
	heap := &GuidHeap{
		Heap{
			useBigIndex: useBigIndex,
			reader:      sr,
		},
	}

	return heap
}

func (heap *GuidHeap) ReadGuid(sr *ShapeReader) Guid {
	var index uint32
	if heap.useBigIndex {
		index = sr.ReadUInt32()
	} else {
		index = uint32(sr.ReadUInt16())
	}
	guid := heap.ReadGuidAtIndex(index)
	return guid
}

func (heap *GuidHeap) ReadGuidAtIndex(index uint32) Guid {
	if err := heap.reader.Seek(int64(index)); err == nil {
		return heap.reader.ReadGuid()
	}

	return ZeroGuid()
}
