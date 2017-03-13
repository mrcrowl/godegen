package cli

type TypeRefRow struct {
	ScopeIndex    ResolutionScopeIndex
	TypeName      string
	TypeNamespace string
}

func (row *TypeRefRow) String() string {
	return row.TypeNamespace + "::" + row.TypeName
}

type ResolutionScopeType uint8

const (
	RSModule ResolutionScopeType = iota
	RSModuleRef
	RSAssemblyRef
	RSTypeRef
)

type ResolutionScopeIndex struct {
	Index uint32
	Type  ResolutionScopeType
}

func NewResolutionScopeIndex(codedIndex CodedIndex) ResolutionScopeIndex {
	return ResolutionScopeIndex{
		Index: codedIndex.Index,
		Type:  ResolutionScopeType(codedIndex.Tag),
	}
}

func readTypeRefRow(
	sr *ShapeReader,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	codedIndex := ReadCodedIndex(sr, tables, TableIdxModule, TableIdxModuleRef, TableIdxAssemblyRef, TableIdxTypeRef)
	scopeIndex := NewResolutionScopeIndex(codedIndex)
	typeName := streams.stringHeap.ReadString(sr)
	typeNamespace := streams.stringHeap.ReadString(sr)
	return &TypeRefRow{
		ScopeIndex:    scopeIndex,
		TypeName:      typeName,
		TypeNamespace: typeNamespace,
	}
}
