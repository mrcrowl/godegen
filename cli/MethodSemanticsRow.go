package cli

type AssociationType uint8

const (
	AssociationEvent AssociationType = iota
	AssociationProperty
)

type AssociationIndex struct {
	Row  uint32
	Type AssociationType
}

type MethodSemanticsRow struct {
	rowNumber       uint32
	Semantics       uint16
	MethodRowNumber uint32
	Association     AssociationIndex
}

func (row *MethodSemanticsRow) String() string {
	return ""
}

func (row *MethodSemanticsRow) RowNumber() uint32 {
	return row.rowNumber
}

func readMethodSemanticsRow(
	sr *ShapeReader,
	rowNumber uint32,
	streams *MetadataStreams,
	tables *TableSet,
) IRow {
	semantics := sr.ReadUInt16()
	methodRowNumber := ReadSimpleIndex(sr, tables, TableIdxMethodDef)
	associationCodexIndex := ReadCodedIndex(sr, tables, TableIdxEvent, TableIdxProperty)
	association := AssociationIndex{associationCodexIndex.Index, AssociationType(associationCodexIndex.Tag)}
	return &MethodSemanticsRow{rowNumber, semantics, methodRowNumber, association}
}
