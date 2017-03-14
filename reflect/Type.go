package reflect

import "godegen/cli"
import "strings"

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
	return typ.Namespace() + "." + typ.Name()
}

type Type interface {
	Name() string
	Namespace() string
	FullName() string
}

type TypeRef struct {
	BaseType
	row *cli.TypeRefRow
}

type TypeDef struct {
	BaseType
	row *cli.TypeDefRow
}

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

func splitFullname(name string) (string, string, bool) {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}

	return "", "", false
}

func newTypeFromDef(typeRow *cli.TypeDefRow, asm *Assembly) Type {
	return &TypeDef{
		BaseType{
			name:      typeRow.TypeName,
			namespace: typeRow.TypeNamespace,
			assembly:  asm,
		},
		typeRow,
	}
}

func newTypeFromRef(typeRow *cli.TypeRefRow, asm *Assembly) Type {
	return &TypeRef{
		BaseType{
			name:      typeRow.TypeName,
			namespace: typeRow.TypeNamespace,
			assembly:  asm,
		},
		typeRow,
	}
}

func (typ *TypeDef) GetMethods() []*Method {
	rows := typ.row.GetMethodRows(typ.assembly.metadata.Tables)
	methods := make([]*Method, len(rows))
	for i, row := range rows {
		methods[i] = newMethod(row, typ.assembly)
	}
	return methods
}
