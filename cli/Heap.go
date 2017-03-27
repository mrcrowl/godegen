package cli

type Heap struct {
	useBigIndex bool
	reader      *ShapeReader
}

func (heap *Heap) GetIndexSizeInBytes() int32 {
	if heap.useBigIndex {
		return 4
	}
	return 2
}

func (heap *Heap) ReadAndDiscard(sr *ShapeReader) {
	if heap.useBigIndex {
		sr.ReadUInt32()
	} else {
		sr.ReadUInt16()
	}
}
