package reflect

type BaseType struct {
	name      string
	namespace string
	assembly  *Assembly
}

func (typ BaseType) Name() string {
	return typ.name
}

func (typ BaseType) Namespace() string {
	return typ.namespace
}

func (typ BaseType) FullName() string {
	// if importedNamespaces[typ.namespace] {
	// 	return typ.name
	// }
	return typ.namespace + "." + typ.name
}

func (typ BaseType) GetFields() []*Field {
	return []*Field{}
}
