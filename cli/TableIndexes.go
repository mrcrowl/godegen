package cli

type CodedIndex struct {
	Index uint32
	Tag   uint8
}

func readSimpleIndex(reader *ShapeReader, tables *TableSet, tableIndex TableIdx) uint32 {
	useBigIndex := tables.useBigIndex(tableIndex)
	var simpleIndex uint32
	if useBigIndex {
		simpleIndex = reader.ReadUInt32()
	} else {
		simpleIndex = uint32(reader.ReadUInt16())
	}
	return simpleIndex
}

func readCodedIndex(reader *ShapeReader, tables *TableSet, tableIndexes ...TableIdx) CodedIndex {
	readInfo := tables.getTableReadInfo(tableIndexes)
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

func readTypeDefOrRefCodedIndex(reader *ShapeReader, tables *TableSet) CodedIndex {
	return readCodedIndex(reader, tables, TableIdxTypeDef, TableIdxTypeRef, TableIdxTypeSpec)
}

func readMemberRefParentCodedIndex(reader *ShapeReader, tables *TableSet) CodedIndex {
	return readCodedIndex(reader, tables, TableIdxMethodDef, TableIdxModuleRef, TableIdxTypeDef, TableIdxTypeRef, TableIdxTypeSpec)
}

func readHasCustomAttributeCodedIndex(reader *ShapeReader, tables *TableSet) CodedIndex {
	return readCodedIndex(reader, tables, TableIdxMethodDef, TableIdxField, TableIdxTypeRef, TableIdxTypeDef, TableIdxParam, TableIdxInterfaceImpl, TableIdxMemberRef, TableIdxModule, TableIdxDeclSecurity, TableIdxProperty, TableIdxEvent, TableIdxStandAloneSig, TableIdxModuleRef, TableIdxTypeSpec, TableIdxAssembly, TableIdxAssemblyRef, TableIdxFile, TableIdxExportedType, TableIdxManifestResource, TableIdxGenericParam, TableIdxGenericParamConstraint, TableIdxMethodSpec)
}

func readCustomAttributeTypeCodedIndex(reader *ShapeReader, tables *TableSet) CodedIndex {
	return readCodedIndex(reader, tables, 0xFF, 0xFF, TableIdxMethodDef, TableIdxMemberRef, 0xFF)
}

func readHasDeclSecurityCodedIndex(reader *ShapeReader, tables *TableSet) CodedIndex {
	return readCodedIndex(reader, tables, TableIdxTypeDef, TableIdxMethodDef, TableIdxAssembly)
}

func readHasFieldMarshalCodedIndex(reader *ShapeReader, tables *TableSet) CodedIndex {
	return readCodedIndex(reader, tables, TableIdxField, TableIdxParam)
}

type TypeOrMethodDefType uint8

const (
	TOMDTypeDef TypeOrMethodDefType = iota
	TOMDMethodDef
)

type TypeOrMethodDefIndex struct {
	Row  uint32
	Type TypeOrMethodDefType
}

func newTypeOrMethodDefIndex(codedIndex CodedIndex) TypeOrMethodDefIndex {
	return TypeOrMethodDefIndex{
		Row:  codedIndex.Index,
		Type: TypeOrMethodDefType(codedIndex.Tag),
	}
}

type TypeDefOrRefType uint8

const (
	TDORTypeDef TypeDefOrRefType = iota
	TDORTypeRef
	TDORTypeSpec
)

type TypeDefOrRefIndex struct {
	Row  uint32
	Type TypeDefOrRefType
}

func newTypeDefOrRefIndex(codedIndex CodedIndex) TypeDefOrRefIndex {
	return TypeDefOrRefIndex{
		Row:  codedIndex.Index,
		Type: TypeDefOrRefType(codedIndex.Tag),
	}
}

type HasConstantType uint8

const (
	HCField HasConstantType = iota
	HCParam
	HCProperty
)

type HasConstantIndex struct {
	Row  uint32
	Type HasConstantType
}

func newHasConstantIndex(codedIndex CodedIndex) HasConstantIndex {
	return HasConstantIndex{
		Row:  codedIndex.Index,
		Type: HasConstantType(codedIndex.Tag),
	}
}
