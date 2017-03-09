package cli

import (
	"math"
)

type StringHeap struct {
	Heap
}

// NewTildeStream =
func NewStringHeap(sr *ShapeReader, useBigIndex bool) *StringHeap {
	heap := &StringHeap{
		Heap{
			useBigIndex: useBigIndex,
			reader:      sr,
		},
	}

	return heap
}

func (heap *StringHeap) ReadString(sr *ShapeReader) string {
	var index uint32
	if heap.useBigIndex {
		index = sr.ReadUInt32()
	} else {
		index = uint32(sr.ReadUInt16())
	}
	str := heap.ReadStringAtIndex(index)
	return str
}

func (heap *StringHeap) ReadStringAtIndex(index uint32) string {
	if err := heap.reader.Seek(int64(index)); err == nil {
		return heap.reader.ReadString(math.MaxInt32)
	}

	return ""
}
