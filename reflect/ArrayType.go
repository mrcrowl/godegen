package reflect

// TODO: add support for multidimensional arrays
type ArrayType struct {
	BaseType
	valueType Type
}

func newArrayType(valueType Type, asm *Assembly) Type {
	return &ArrayType{
		BaseType{
			name:      valueType.Name(),
			namespace: valueType.Namespace(),
			assembly:  asm,
		},
		valueType,
	}
}

func (array *ArrayType) rowNumber() uint32 {
	return 0
}

func (array *ArrayType) Name() string {
	valueTypeName := array.valueType.Name()
	return valueTypeName + "[]"
}

func (array *ArrayType) FullName() string {
	return array.namespace + "." + array.Name()
}

func (array *ArrayType) GetMethods() []*Method {
	return []*Method{}
}

func (array *ArrayType) GetFields() []*Field {
	return []*Field{}
}
