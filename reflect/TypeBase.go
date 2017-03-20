package reflect

type TypeBase struct {
	name      string
	namespace string
	assembly  *Assembly
}

func (typ TypeBase) Name() string {
	return typ.name
}

func (typ TypeBase) Namespace() string {
	return typ.namespace
}

func (typ TypeBase) FullName() string {
	// if importedNamespaces[typ.namespace] {
	// 	return typ.name
	// }
	return typ.namespace + "." + typ.name
}

func (typ TypeBase) GetFields() []*Field {
	return []*Field{}
}
