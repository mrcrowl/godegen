package reflect

import (
	"godegen/cli"
)

type Property struct {
	name      string
	signature *PropertySig
	static    bool
	initOnly  bool
	literal   bool
	special   bool
	getter    *PropertyMethod
	setter    *PropertyMethod
}

type PropertyMethodKind uint8

const (
	PMKNone   PropertyMethodKind = 0
	PMKSetter PropertyMethodKind = 1
	PMKGetter PropertyMethodKind = 2
	PMKMask   uint16             = uint16(PMKGetter | PMKSetter)
)

type PropertyMethod struct {
	Method
	kind PropertyMethodKind
}

const (
	PropertyAttributesStatic      uint16 = 0x10
	PropertyAttributesInitOnly    uint16 = 0x20
	PropertyAttributesLiteral     uint16 = 0x40
	PropertyAttributesSpecialName uint16 = 0x200
	PropertyAttributesHasDefault  uint16 = 0x1000
)

func newProperty(propRow *cli.PropertyRow, semanticRowsByProp map[uint32][]*cli.MethodSemanticsRow, asm *Assembly) *Property {
	sigBlob := propRow.GetSignatureBlob()
	sigReader := NewSignatureReader(sigBlob, asm)
	static := (propRow.Flags & PropertyAttributesStatic) > 0
	initOnly := (propRow.Flags & PropertyAttributesInitOnly) > 0
	literal := (propRow.Flags & PropertyAttributesLiteral) > 0
	special := (propRow.Flags & PropertyAttributesSpecialName) > 0

	methodSemanticsForProp := semanticRowsByProp[propRow.RowNumber()]
	methodDefTable := asm.metadata.Tables.GetTable(cli.TableIdxMethodDef)

	var getter *PropertyMethod
	var setter *PropertyMethod
	for _, semantics := range methodSemanticsForProp {
		methodDefRow := methodDefTable.GetRow(semantics.MethodRowNumber).(*cli.MethodDefRow)
		method := *newMethod(methodDefRow, asm)

		kind := PropertyMethodKind(semantics.Semantics & PMKMask)
		switch kind {
		case PMKGetter:
			getter = &PropertyMethod{method, PMKGetter}
		case PMKSetter:
			setter = &PropertyMethod{method, PMKSetter}
		}
	}

	return &Property{
		name:      propRow.Name,
		signature: sigReader.ReadPropertySignature(),
		static:    static,
		initOnly:  initOnly,
		literal:   literal,
		getter:    getter,
		special:   special,
		setter:    setter,
	}
}

func (prop *Property) Name() string {
	return prop.name
}

func (prop *Property) Type() Type {
	return prop.signature.propertyType
}

func (prop *Property) Getter() *PropertyMethod {
	return prop.getter
}

func (prop *Property) HasPublicGetter() bool {
	if prop.getter == nil {
		return false
	}

	return prop.getter.memberAccess == Public
}

func (prop *Property) Setter() *PropertyMethod {
	return prop.setter
}
