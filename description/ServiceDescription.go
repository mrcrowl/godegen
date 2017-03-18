package description

type ServiceDescription struct {
	Namespaces []*Namespace
}

type Namespace struct {
	Name       string
	FullName   string
	Namespaces []*Namespace
	Services   []*Service
}

type Service struct {
	Name      string
	Namespace string
	FullName  string
	Methods   []*Method
}

type Method struct {
	Name          string
	TypeName      string
	TypeNamespace string
	Args          []*Arg
}

type Arg struct {
	Name                 string
	TypeName             string
	TypeNamespace        string
	ElementTypeName      string
	ElementTypeNamespace string
}

type DataType struct {
	Name                 string
	Namespace            string
	FullName             string
	BaseTypeName         string
	BaseTypeNamespace    string
	ElementTypeName      string
	ElementTypeNamespace string
}

type Field struct {
	Name       string
	Type       string
	SourceType string
}

type Const struct {
	Name       string
	Type       string
	SourceType string
	Value      string
}
