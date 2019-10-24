package reflect

import (
	"educationperfect.com/godegen/cli"
	"strings"
)

const CONSTRUCTOR_NAME = ".ctor"

type NormalType struct {
	TypeBase
	extends Type
	row     *cli.TypeDefRow
}

func (def *NormalType) RowNumber() uint32 {
	return def.row.RowNumber()
}

func (def *NormalType) Base() Type {
	return def.extends
}

func newTypeFromDef(typeRow *cli.TypeDefRow, extendsType Type, asm *Assembly) Type {
	return &NormalType{
		TypeBase{
			name:      typeRow.TypeName,
			namespace: typeRow.TypeNamespace,
			assembly:  asm,
		},
		extendsType,
		typeRow,
	}
}

func (typ *NormalType) GetMethods() []*Method {
	rows := typ.row.GetMethodRows(typ.assembly.metadata.Tables)
	methods := make([]*Method, 0, len(rows))
	count := 0
	for _, row := range rows {
		method := newMethod(row, typ.assembly)
		includeMethod := method.memberAccess == Public && method.Name() != CONSTRUCTOR_NAME && !method.static
		if includeMethod {
			methods = append(methods, method)
			count++
		}
	}
	return methods[:count]
}

func (typ *NormalType) GetFields() []*Field {
	return typ.GetFieldsWithOptions(false, true, false)
}

func (typ *NormalType) GetFieldsWithOptions(includeNonPublic bool, includeInstance bool, includeStatic bool) []*Field {
	rows := typ.row.GetFieldRows(typ.assembly.metadata.Tables)
	fields := make([]*Field, 0, len(rows))
	count := 0
	for _, row := range rows {
		field := newField(row, typ.assembly)
		includeField := ((includeNonPublic || field.memberAccess == Public) &&
			(includeInstance || field.static) &&
			(includeStatic || !field.static))
		if includeField {
			fields = append(fields, field)
			count++
		}
	}
	return fields[:count]
}

func (typ *NormalType) GetProperties() []*Property {
	assembly := typ.assembly
	rows := typ.row.GetPropertyRows(assembly.metadata.Tables)
	numRows := len(rows)
	if numRows > 0 {
		propRowRange := cli.NewRowRange(rows[0].RowNumber(), rows[numRows-1].RowNumber()+1)
		semanticRows := typ.assembly.metadata.GetMethodSemanticsRowsForProps(propRowRange)

		properties := make([]*Property, 0, len(rows))
		count := 0
		for _, row := range rows {
			property := newProperty(row, semanticRows, typ.assembly)
			if property.HasPublicGetter() {
				properties = append(properties, property)
				count++
			}
		}
		return properties[:count]
	}

	return []*Property{}
}

func splitFullname(fullname string) (string, string, bool) {
	dot := strings.LastIndex(fullname, ".")
	if dot >= 0 {
		namespace := fullname[0:dot]
		name := fullname[dot+1:]
		return namespace, name, true
	}

	return "", "", false
}
