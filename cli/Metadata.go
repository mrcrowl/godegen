package cli

// Metadata =
type Metadata struct {
	Signature     uint32
	MajorVersion  uint16
	MinorVersion  uint16
	Reserved      uint32
	VersionLength uint32
	Version       string
	Tables        *TableSet

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

// NewMetadata =
func NewMetadata(textSection *TextSection, metadataRVA RVA) *Metadata {
	metaDataReader := textSection.GetReaderAt(metadataRVA.VirtualAddress)
	mr := NewShapeReader(metaDataReader)

	md := new(Metadata)
	md.baseAddress = metadataRVA.VirtualAddress
	md.Signature = mr.ReadUInt32()
	md.MajorVersion = mr.ReadUInt16()
	md.MinorVersion = mr.ReadUInt16()
	md.Reserved = mr.ReadUInt32()
	md.VersionLength = mr.ReadUInt32()
	md.Version = mr.ReadUTF8(md.VersionLength)
	skipOffset1 := calcSkipBytes(md.VersionLength, 4)
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
	}
	md.Tables = NewTableSet(&streams)
	md.Tables.ReadAll(&streams)

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
