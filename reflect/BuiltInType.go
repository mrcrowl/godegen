package reflect

type BuiltInType struct {
	TypeBase
	shortName string
}

func newBuiltInType(fullname string, shortName string) Type {
	if namespace, name, ok := splitFullname(fullname); ok {
		return &BuiltInType{
			TypeBase{
				name,
				namespace,
				nil,
			},
			shortName,
		}
	}
	return nil
}

func newBuiltInTypeWithoutAlias(fullname string) Type {
	if namespace, name, ok := splitFullname(fullname); ok {
		return &BuiltInType{
			TypeBase{
				name,
				namespace,
				nil,
			},
			"",
		}
	}
	return nil
}

func (bi *BuiltInType) FullName() string {
	// if len(bi.shortName) > 0 {
	// 	return bi.shortName
	// }
	return bi.TypeBase.FullName()
}

func (bi *BuiltInType) RowNumber() uint32 {
	return 0
}

func (bi *BuiltInType) Base() Type {
	return nil
}

func (bi *BuiltInType) GetMethods() []*Method {
	return []*Method{}
}

func (bi *BuiltInType) GetProperties() []*Property {
	return []*Property{}
}
