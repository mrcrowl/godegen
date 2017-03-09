package cli

type CodedIndex struct {
	Index uint32
	Tag   uint8
}

func ReadCodedIndex(reader *ShapeReader, tables *TableSet, tableIndexes ...TableIdx) CodedIndex {
	readInfo := tables.GetTableReadInfo(tableIndexes)
	var codedIndex uint32
	if readInfo.UseBigIndex {
		codedIndex = reader.ReadUInt32()
	} else {
		codedIndex = uint32(reader.ReadUInt16())
	}

	tagMask := uint32((1 << readInfo.NumTagBits) - 1)
	tag := uint8(codedIndex & tagMask)
	index := codedIndex >> readInfo.NumTagBits
	return CodedIndex{Index: index, Tag: tag}
}
