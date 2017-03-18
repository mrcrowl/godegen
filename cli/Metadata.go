package cli

// Metadata =
type Metadata struct {
	signature    uint32
	majorVersion uint16
	minorVersion uint16
	reserved     uint32
	version      string
	Tables       *TableSet

	baseAddress   uint32
	streamHeaders []streamHeader
}

const smallHeapSize = 0x10000

type streamHeader struct {
	offset uint32
	size   uint32
	name   string
}

type MetadataStreams struct {
	tildeStream *TildeStream
	stringHeap  *StringHeap
	guidHeap    *GuidHeap
	blobHeap    *BlobHeap
}

func calcSkipBytes(n uint32, multiple uint32) int32 {
	if multiple == 0 {
		return int32(n)
	}

	remainder := n % multiple
	if remainder == 0 {
		return 0
	}

	return int32(multiple - remainder)
}

func newMetadata(textSection *TextSection, metadataRVA RVA) *Metadata {
	metaDataReader := textSection.GetReaderAt(metadataRVA.VirtualAddress)
	mr := NewShapeReader(metaDataReader)

	md := new(Metadata)
	md.baseAddress = metadataRVA.VirtualAddress
	md.signature = mr.ReadUInt32()
	md.majorVersion = mr.ReadUInt16()
	md.minorVersion = mr.ReadUInt16()
	md.reserved = mr.ReadUInt32()

	versionLength := mr.ReadUInt32()
	md.version = mr.ReadUTF8(versionLength)
	skipOffset1 := calcSkipBytes(versionLength, 4)
	mr.Skip(skipOffset1)
	mr.Skip(2)

	// stream
	numStreams := uint32(mr.ReadUInt16())
	md.streamHeaders = readStreamHeaders(mr, numStreams)

	// get known streams
	streams := MetadataStreams{
		tildeStream: md.getTildeStream(textSection),
		stringHeap:  md.getStringHeap(textSection),
		guidHeap:    md.getGuidHeap(textSection),
		blobHeap:    md.getBlobHeap(textSection),
	}
	md.Tables = NewTableSet(&streams)
	md.Tables.readAll(&streams)

	return md
}

func readStreamHeaders(mr *ShapeReader, numStreams uint32) []streamHeader {
	streamHeaders := make([]streamHeader, numStreams)
	for i := uint32(0); i < numStreams; i++ {
		offset := mr.ReadUInt32()
		size := mr.ReadUInt32()
		name := mr.ReadString(32)
		streamHeaders[i] = streamHeader{name: name, offset: offset, size: size}

		// skip to next stream
		nameLength := len(name) + 1
		skipOffset2 := calcSkipBytes(uint32(nameLength), 4)
		mr.Skip(skipOffset2)
	}
	return streamHeaders
}

func (md *Metadata) getStreamHeader(name string) (streamHeader, bool) {
	for _, header := range md.streamHeaders {
		if header.name == name {
			return header, true
		}
	}
	return streamHeader{}, false
}

func (md *Metadata) getTildeStream(textSection *TextSection) *TildeStream {
	if header, ok := md.getStreamHeader("#~"); ok {
		streamReader := textSection.GetReaderAt(md.baseAddress + header.offset)
		sr := NewShapeReader(streamReader)
		return NewTildeStream(sr)
	}

	return nil
}

func (md *Metadata) getStringHeap(textSection *TextSection) *StringHeap {
	if header, ok := md.getStreamHeader("#Strings"); ok {
		heapReader := textSection.GetReaderAt(md.baseAddress + header.offset)
		sr := NewShapeReader(heapReader)
		useBigIndex := header.size >= smallHeapSize
		return NewStringHeap(sr, useBigIndex)
	}

	return nil
}

func (md *Metadata) getGuidHeap(textSection *TextSection) *GuidHeap {
	if header, ok := md.getStreamHeader("#GUID"); ok {
		heapReader := textSection.GetReaderAt(md.baseAddress + header.offset)
		sr := NewShapeReader(heapReader)
		useBigIndex := header.size >= smallHeapSize
		return NewGuidHeap(sr, useBigIndex)
	}

	return nil
}

func (md *Metadata) getBlobHeap(textSection *TextSection) *BlobHeap {
	if header, ok := md.getStreamHeader("#Blob"); ok {
		heapReader := textSection.GetReaderAt(md.baseAddress + header.offset)
		sr := NewShapeReader(heapReader)
		useBigIndex := header.size >= smallHeapSize
		return NewBlobHeap(sr, useBigIndex)
	}

	return nil
}

func (md *Metadata) GetMethodSemanticsRowsForProps(propRowRange RowRange) map[uint32][]*MethodSemanticsRow {
	rowsByProp := make(map[uint32][]*MethodSemanticsRow, propRowRange.count())
	for _, row := range md.Tables.GetTable(TableIdxMethodSemantics).rows {
		methodSemRow := row.(*MethodSemanticsRow)
		association := methodSemRow.Association
		if association.Type == AssociationProperty && association.Row >= propRowRange.from && association.Row < propRowRange.to {
			var existingMethodSemRows []*MethodSemanticsRow
			if existingMethodSemRows = rowsByProp[association.Row]; existingMethodSemRows == nil {
				existingMethodSemRows = make([]*MethodSemanticsRow, 0, 2)
				rowsByProp[association.Row] = append(existingMethodSemRows, methodSemRow)
			} else {
				rowsByProp[association.Row] = append(existingMethodSemRows, methodSemRow)
			}
		}
	}
	return rowsByProp
}
