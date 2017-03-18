package reflect

import (
	"godegen/cli"
	"strings"
)

const CONSTRUCTOR_NAME = ".ctor"

type NormalType struct {
	BaseType
	row *cli.TypeDefRow
}

func (def *NormalType) rowNumber() uint32 {
	return def.row.RowNumber()
}

func newTypeFromDef(typeRow *cli.TypeDefRow, asm *Assembly) Type {
	return &NormalType{
		BaseType{
			name:      typeRow.TypeName,
			namespace: typeRow.TypeNamespace,
			assembly:  asm,
		},
		typeRow,
	}
}

func (typ *NormalType) GetMethods() []*Method {
	rows := typ.row.GetMethodRows(typ.assembly.metadata.Tables)
	methods := make([]*Method, 0, len(rows))
	count := 0
	for _, row := range rows {
		method := newMethod(row, typ.assembly)
		includeMethod := method.memberAccess == Public && method.Name() != CONSTRUCTOR_NAME
		if includeMethod {
			methods = append(methods, method)
			count++
		}
	}
	return methods[:count]
}

func (typ *NormalType) GetFields() []*Field {
	rows := typ.row.GetFieldRows(typ.assembly.metadata.Tables)
	fields := make([]*Field, 0, len(rows))
	count := 0
	for _, row := range rows {
		field := newField(row, typ.assembly)
		if field.memberAccess == Public {
			fields = append(fields, field)
			count++
		}
	}
	return fields[:count]
}

func splitFullname(name string) (string, string, bool) {
	parts := strings.SplitN(name, ".", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}

	return "", "", false
}
