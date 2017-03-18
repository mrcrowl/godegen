package reflect

import (
	"bytes"
	"strings"
)

type GenericType struct {
	BaseType
	templateType Type
	numArgs      uint32
	argTypes     []Type
}

func newGenericType(templateType Type, argTypes []Type, asm *Assembly) Type {
	return &GenericType{
		BaseType{
			name:      templateType.Name(),
			namespace: templateType.Namespace(),
			assembly:  asm,
		},
		templateType,
		uint32(len(argTypes)),
		argTypes,
	}
}

func (gen *GenericType) Name() string {
	var buffer bytes.Buffer
	lexicalName := gen.LexicalName()
	buffer.WriteString(lexicalName)
	buffer.WriteByte('<')
	for i, arg := range gen.argTypes {
		buffer.WriteString(arg.FullName())
		if uint32(i+1) < gen.numArgs {
			buffer.WriteByte(',')
		}
	}
	buffer.WriteByte('>')
	return buffer.String()
}

func (gen *GenericType) CommonName() string {
	lexicalName := gen.LexicalName()
	if lexicalName == "Nullable" && gen.namespace == "System" && gen.numArgs == 1 {
		argName := gen.argTypes[0].FullName()
		return argName + "?"
	}

	return gen.Name()
}

func (gen *GenericType) LexicalName() string {
	return strings.SplitN(gen.name, "`", 2)[0]
}

func (gen *GenericType) FullName() string {
	// if importedNamespaces[gen.namespace] {
	// 	return gen.Name()
	// }
	return gen.namespace + "." + gen.Name()
}

func (gen *GenericType) ArgumentTypes() []Type {
	return gen.argTypes
}

func (gen *GenericType) RowNumber() uint32 {
	return 0
}

func (gen *GenericType) GetMethods() []*Method {
	return gen.templateType.GetMethods()
}

func (gen *GenericType) GetFields() []*Field {
	return gen.templateType.GetFields()
}

func (gen *GenericType) GetProperties() []*Property {
	return gen.templateType.GetProperties()
}
