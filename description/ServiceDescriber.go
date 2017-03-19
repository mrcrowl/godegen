package description

import (
	"godegen/reflect"
	"strings"
)

type ServiceDescriber struct {
	assemblyFile    *reflect.Assembly
	typeMapper      func(reflect.Type) string
	namespaceMapper func(string) string
}

func NewServiceDescriber(assemblyFilepath string, assemblyName string, typeMapper func(reflect.Type) string, namespaceMapper func(string) string) *ServiceDescriber {
	loader := reflect.NewAssemblyLoader(assemblyFilepath)
	assemblyFile, _ := loader.Load(assemblyName)
	return &ServiceDescriber{assemblyFile, typeMapper, namespaceMapper}
}

func (res *ServiceDescriber) Describe(serviceTypeName string) (*ServiceDescription, error) {
	serviceType := res.assemblyFile.GetType(serviceTypeName)
	resolvedTypes := ResolveServiceDependencyTypes(res.assemblyFile, serviceType)
	// printTypes(resolvedTypes)

	description := res.createDescriptionOfTypes(resolvedTypes, serviceType)
	return description, nil
}

func (res *ServiceDescriber) createDescriptionOfTypes(types []reflect.Type, serviceType reflect.Type) *ServiceDescription {
	rootNamespaces, namespaceMap := res.buildNamespaceTree(types, serviceType)

	// add service
	serviceTypeNamespace := res.mapNamespace(serviceType.Namespace())
	serviceNamespace := namespaceMap[serviceTypeNamespace]
	service := res.createService(serviceType)
	serviceNamespace.addService(service)

	// fmt.Println(rootNamespaces)
	return &ServiceDescription{
		Namespaces: rootNamespaces,
	}
}

func (res *ServiceDescriber) createService(serviceType reflect.Type) *Service {
	methods := res.collectTypeMethods(serviceType)

	return &Service{
		*res.createDataType(serviceType),
		methods,
	}
}

func (res *ServiceDescriber) buildNamespaceTree(types []reflect.Type, serviceType reflect.Type) ([]*Namespace, map[string]*Namespace) {
	var namespaceSeen = map[string]*Namespace{}
	var distinctNamespaces []*Namespace
	var rootNamespaces []*Namespace

	for _, typ := range append(types, serviceType) {
		var childNamespace *Namespace
		var namespace *Namespace
		var found bool
		var nsName = res.mapNamespace(typ.Namespace())
		var nsOriginal = nsName
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

		namespace = namespaceSeen[nsOriginal]

		// add type to namespace
		if typ != serviceType {
			dataType := res.createDataType(typ)
			namespace.DataTypes = append(namespace.DataTypes, dataType)
		}
	}

	return rootNamespaces, namespaceSeen
}

func (res *ServiceDescriber) createDataType(typ reflect.Type) *DataType {
	fields := res.collectTypeFields(typ)

	return &DataType{
		DataTypeReference{
			Name:      typ.Name(),
			Namespace: res.mapNamespace(typ.Namespace()),
			// QualifiedName: res.mapNamespace(typ.FullName()),
		},
		nil, // TODO: get base type
		fields,
	}
}

// func createDataTypeReference(typ reflect.Type) *DataTypeReference {
// 	var elementDataType *DataTypeReference
// 	if elementType, isCollection := isCollectionType(typ); isCollection {
// 		elementDataType = createDataTypeReference(elementType)
// 	}

// 	return &DataTypeReference{
// 		Name:          typ.Name(),
// 		Namespace:     typ.Namespace(),
// 		QualifiedName: typ.FullName(),
// 		ElementType:   elementDataType,
// 	}
// }

func (res *ServiceDescriber) collectTypeFields(typ reflect.Type) []*Field {
	var fields []*Field

	for _, field := range typ.GetFields() {
		fields = append(fields, res.createFieldFromField(field))
	}

	for _, property := range typ.GetProperties() {
		fields = append(fields, res.createFieldFromProperty(property))
	}

	return fields
}

func (res *ServiceDescriber) mapType(typ reflect.Type) string {
	return res.typeMapper(typ)
}

func (res *ServiceDescriber) mapNamespace(namespace string) string {
	return res.namespaceMapper(namespace)
}

func (res *ServiceDescriber) collectTypeMethods(typ reflect.Type) []*Method {
	var methods []*Method

	for _, method := range typ.GetMethods() {
		args := res.collectMethodArgs(method)

		meth := &Method{
			Name: method.Name(),
			Type: res.mapType(method.ReturnType()),
			Args: args,
		}
		methods = append(methods, meth)
	}

	return methods
}

func (res *ServiceDescriber) collectMethodArgs(method *reflect.Method) []*Arg {
	var args []*Arg

	for _, param := range method.Parameters() {
		arg := &Arg{
			Name: param.Name(),
			Type: res.mapType(param.Type()),
		}
		args = append(args, arg)
	}

	return args
}

func (res *ServiceDescriber) createFieldFromField(field *reflect.Field) *Field {
	return &Field{
		Name: field.Name(),
		Type: res.mapType(field.Type()),
	}
}

func (res *ServiceDescriber) createFieldFromProperty(property *reflect.Property) *Field {
	return &Field{
		Name: property.Name(),
		Type: res.mapType(property.Type()),
	}
}

func newNamespace(namespace string) *Namespace {
	name := getLastNamespaceSegment(namespace)
	return &Namespace{
		Name:          name,
		qualifiedName: namespace,
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

// func printTypes(resolvedTypes []reflect.Type) {
// 	for _, t := range resolvedTypes {
// 		fmt.Println(t.FullName() + ":")
// 		for _, field := range t.GetFields() {
// 			fieldType := field.Type()
// 			if collectionType, isCollection := isCollectionType(fieldType); isCollection {
// 				fmt.Println("\t" + field.Name() + ": " + collectionType.FullName() + "[]")
// 			} else {
// 				fmt.Println("\t" + field.Name() + ": " + field.Type().FullName())
// 			}
// 		}
// 	}
// }
