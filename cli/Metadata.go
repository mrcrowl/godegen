package cli

// Metadata =
type Metadata struct {
	Signature     uint32
	MajorVersion  uint16
	MinorVersion  uint16
	Reserved      uint32
	VersionLength uint32
	Version       string
	NumStreams    uint16
}

func (md *Metadata) NumberOfStreams() uint16 {
	return 0
}

// NewMetadata =
func NewMetadata(textSection *TextSection, metadataRVA RVA) *Metadata {
	metaDataReader := textSection.GetReaderAt(metadataRVA.VirtualAddress)
	mr := NewMetaReader(metaDataReader)

	md := new(Metadata)
	md.Signature = mr.ReadUInt32()
	md.MajorVersion = mr.ReadUInt16()
	md.MinorVersion = mr.ReadUInt16()
	md.Reserved = mr.ReadUInt32()
	md.VersionLength = mr.ReadUInt32()
	md.Version = mr.ReadUTF8(md.VersionLength)
	alignmentSize := md.VersionLength / 4 * 4
	offset := int32(alignmentSize - md.VersionLength)
	mr.Skip(offset)
	mr.Skip(2)
	md.NumStreams = mr.ReadUInt16()
	return md
}
