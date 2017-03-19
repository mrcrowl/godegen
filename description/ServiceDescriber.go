package description

import (
	"fmt"
	"godegen/reflect"
	"strings"
)

type ServiceDescriber struct {
	assemblyFile *reflect.Assembly
}

func NewServiceDescriber(assemblyFilepath string, assemblyName string) *ServiceDescriber {
	loader := reflect.NewAssemblyLoader(assemblyFilepath)
	assemblyFile, _ := loader.Load(assemblyName)
	return &ServiceDescriber{assemblyFile}
}

func (res *ServiceDescriber) Describe(serviceTypeName string) (*ServiceDescription, error) {
	serviceType := res.assemblyFile.GetType(serviceTypeName)
	resolvedTypes := ResolveServiceDependencyTypes(res.assemblyFile, serviceType)
	// printTypes(resolvedTypes)

	description := res.createDescriptionOfTypes(resolvedTypes, serviceType)
	return description, nil
}

func (res *ServiceDescriber) createDescriptionOfTypes(types []reflect.Type, serviceType reflect.Type) *ServiceDescription {
	rootNamespaces, namespaceMap := buildNamespaceTree(types, serviceType)

	// add service
	serviceNamespace := namespaceMap[serviceType.Namespace()]
	service := createService(serviceType)
	serviceNamespace.addService(service)

	// fmt.Println(rootNamespaces)
	return &ServiceDescription{
		Namespaces: rootNamespaces,
	}
}

func createService(serviceType reflect.Type) *Service {
	methods := collectTypeMethods(serviceType)

	return &Service{
		*createDataTypeReference(serviceType),
		methods,
	}
}

func buildNamespaceTree(types []reflect.Type, serviceType reflect.Type) ([]*Namespace, map[string]*Namespace) {
	var namespaceSeen = map[string]*Namespace{}
	var distinctNamespaces []*Namespace
	var rootNamespaces []*Namespace

	for _, typ := range append(types, serviceType) {
		var childNamespace *Namespace
		var namespace *Namespace
		var found bool
		nsName := typ.Namespace()
		for {
			if namespace, found = namespaceSeen[nsName]; !found {
				namespace = newNamespace(nsName)
			}

			if childNamespace != nil {
				namespace.addChild(childNamespace)
			}

			if found {
				break
			}

			namespaceSeen[nsName] = namespace
			distinctNamespaces = append(distinctNamespaces, namespace)
			if namespace.isRoot() {
				rootNamespaces = append(rootNamespaces, namespace)
				break
			}

			nsName = getParentNamespace(nsName)
			childNamespace = namespace
		}

		namespace = namespaceSeen[nsName]

		// add type to namespace
		if typ != serviceType {
			dataType := createDataType(typ)
			namespace.DataTypes = append(namespace.DataTypes, dataType)
		}
	}

	return rootNamespaces, namespaceSeen
}

func createDataType(typ reflect.Type) *DataType {
	fields := collectTypeFields(typ)

	return &DataType{
		*createDataTypeReference(typ),
		nil, // TODO: get base type
		fields,
	}
}

func createDataTypeReference(typ reflect.Type) *DataTypeReference {
	var elementDataType *DataTypeReference
	if elementType, isCollection := isCollectionType(typ); isCollection {
		elementDataType = createDataTypeReference(elementType)
	}

	return &DataTypeReference{
		Name:          typ.Name(),
		Namespace:     typ.Namespace(),
		QualifiedName: typ.FullName(),
		ElementType:   elementDataType,
	}
}

func collectTypeFields(typ reflect.Type) []*Field {
	var fields []*Field

	for _, field := range typ.GetFields() {
		fields = append(fields, createFieldFromField(field))
	}

	for _, property := range typ.GetProperties() {
		fields = append(fields, createFieldFromProperty(property))
	}

	return fields
}

func collectTypeMethods(typ reflect.Type) []*Method {
	var methods []*Method

	for _, method := range typ.GetMethods() {
		args := collectMethodArgs(method)

		meth := &Method{
			Name: method.Name(),
			Type: createDataTypeReference(method.ReturnType()),
			Args: args,
		}
		methods = append(methods, meth)
	}

	return methods
}

func collectMethodArgs(method *reflect.Method) []*Arg {
	var args []*Arg

	for _, param := range method.Parameters() {
		arg := &Arg{
			Name: param.Name(),
			Type: createDataTypeReference(param.Type()),
		}
		args = append(args, arg)
	}

	return args
}

func createFieldFromField(field *reflect.Field) *Field {
	return &Field{
		Name: field.Name(),
		Type: createDataTypeReference(field.Type()),
	}
}

func createFieldFromProperty(property *reflect.Property) *Field {
	return &Field{
		Name: property.Name(),
		Type: createDataTypeReference(property.Type()),
	}
}

func newNamespace(namespace string) *Namespace {
	name := getLastNamespaceSegment(namespace)
	return &Namespace{
		Name:          name,
		QualifiedName: namespace,
	}
}

func getLastNamespaceSegment(namespace string) string {
	dot := strings.LastIndex(namespace, ".")
	if dot >= 0 {
		return namespace[dot+1:]
	}
	return namespace
}

func getParentNamespace(namespace string) string {
	dot := strings.LastIndex(namespace, ".")
	if dot < 0 {
		return ""
	}
	return namespace[:dot]
}

func printTypes(resolvedTypes []reflect.Type) {
	for _, t := range resolvedTypes {
		fmt.Println(t.FullName() + ":")
		for _, field := range t.GetFields() {
			fieldType := field.Type()
			if collectionType, isCollection := isCollectionType(fieldType); isCollection {
				fmt.Println("\t" + field.Name() + ": " + collectionType.FullName() + "[]")
			} else {
				fmt.Println("\t" + field.Name() + ": " + field.Type().FullName())
			}
		}
	}
}

func isCollectionType(typ reflect.Type) (reflect.Type, bool) {
	if array, isArray := typ.(*reflect.ArrayType); isArray {
		return array.ValueType(), true
	}

	if generic, isGeneric := typ.(*reflect.GenericType); isGeneric {
		if generic.BaseType.FullName() == "System.Collections.Generic.List`1" {
			return generic.ArgumentTypes()[0], true
		}
	}

	return nil, false
}
