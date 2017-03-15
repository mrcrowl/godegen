package reflect

import "godegen/cli"
import "strings"
import "bytes"

var importedNamespaces = map[string]bool{
	"System":                     true,
	"System.Threading.Tasks":     true,
	"System.Collections.Generic": true,
}

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
	if importedNamespaces[typ.namespace] {
		return typ.name
	}
	return typ.namespace + "." + typ.name
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

func (bi *BuiltInType) FullName() string {
	if len(bi.shortName) > 0 {
		return bi.shortName
	}
	return bi.BaseType.FullName()
}

type GenericType struct {
	BaseType
	numArgs  int
	argTypes []Type
}

func (gen *GenericType) Name() string {
	var buffer bytes.Buffer
	lexicalName := gen.LexicalName()
	if lexicalName == "Nullable" && gen.namespace == "System" && gen.numArgs == 1 {
		argName := gen.argTypes[0].FullName()
		return argName + "?"
	}
	buffer.WriteString(lexicalName)
	buffer.WriteByte('<')
	for i, arg := range gen.argTypes {
		buffer.WriteString(arg.FullName())
		if (i + 1) < gen.numArgs {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte('>')
	return buffer.String()
}

func (gen *GenericType) LexicalName() string {
	return strings.SplitN(gen.name, "`", 2)[0]
}

func (gen *GenericType) FullName() string {
	if importedNamespaces[gen.namespace] {
		return gen.Name()
	}
	return gen.namespace + "." + gen.Name()
}

func newGenericType(name string, namespace string, argTypes []Type, asm *Assembly) Type {
	return &GenericType{
		BaseType{
			name:      name,
			namespace: namespace,
			assembly:  asm,
		},
		len(argTypes),
		argTypes,
	}
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
	methods := make([]*Method, 0, len(rows))
	count := 0
	for _, row := range rows {
		method := newMethod(row, typ.assembly)
		if method.memberAccess == Public {
			methods = append(methods, method)
			count++
		}
	}
	return methods[:count]
}

func splitFullname(name string) (string, string, bool) {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}

	return "", "", false
}
