package reflect

type typeCache struct {
	builtInTypes []Type
	typesByName  map[string]Type
}

func newTypeCache() *typeCache {
	return &typeCache{
		typesByName:  make(map[string]Type),
		builtInTypes: getBuiltInTypes(),
	}
}

func (cache *typeCache) getBuiltIn(id byte) Type {
	if id >= 0 && id <= 0x1c {
		return cache.builtInTypes[id]
	}
	return nil
}

func (cache *typeCache) get(name string) Type {
	if value, ok := cache.typesByName[name]; ok {
		return value
	}
	return nil
}

func (cache *typeCache) set(name string, t Type) {
	cache.typesByName[name] = t
}

func getBuiltInTypes() []Type {
	return []Type{
		nil, // 0x00
		newBuiltInType("System.Void", "void"),     // 0x01
		newBuiltInType("System.Boolean", "bool"),  // 0x02
		newBuiltInType("System.Char", "char"),     //0x03
		newBuiltInType("System.SByte", "sbyte"),   //0x04
		newBuiltInType("System.Byte", "byte"),     //0x05
		newBuiltInType("System.Int16", "short"),   //0x06
		newBuiltInType("System.UInt16", "ushort"), //0x07
		newBuiltInType("System.Int32", "int"),     //0x08
		newBuiltInType("System.UInt32", "uint"),   //0x09
		newBuiltInType("System.Int64", "long"),    //0x0a
		newBuiltInType("System.UInt64", "ulong"),  //0x0b
		newBuiltInType("System.Single", "float"),  //0x0c
		newBuiltInType("System.Double", "double"), //0x0d
		newBuiltInType("System.String", "string"), //0x0e
		nil, // 0x0f
		nil, // 0x10
		nil, // 0x11
		nil, // 0x12
		nil, // 0x13
		nil, // 0x14
		nil, // 0x15
		newBuiltInTypeWithoutAlias("System.TypedReference"), //0x16
		nil, // 0x17
		newBuiltInTypeWithoutAlias("System.IntPtr"),  //0x18
		newBuiltInTypeWithoutAlias("System.UIntPtr"), //0x19
		nil, // 0x1a
		nil, // 0x1b
		newBuiltInType("System.Object", "object"), //0x1c
	}
}
