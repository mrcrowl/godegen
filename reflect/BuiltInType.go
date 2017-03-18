package reflect

type BuiltInType struct {
	BaseType
	shortName string
}

func newBuiltInType(fullname string, shortName string) Type {
	if namespace, name, ok := splitFullname(fullname); ok {
		return &BuiltInType{
			BaseType{
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
			BaseType{
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
	if len(bi.shortName) > 0 {
		return bi.shortName
	}
	return bi.BaseType.FullName()
}

func (bi *BuiltInType) rowNumber() uint32 {
	return 0
}

func (bi *BuiltInType) GetMethods() []*Method {
	return []*Method{}
}
