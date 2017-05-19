package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ServiceDescription .
type ServiceDescription struct {
	Namespaces []*Namespace
}

// JSON outputs to JSON format
func (desc *ServiceDescription) JSON() string {
	// b, _ := json.MarshalIndent(desc, "", "\t")
	// b, _ := json.Marshal(desc)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(desc); err != nil {
		fmt.Println(err)
	}
	return buf.String()
}

// Namespace is a type
type Namespace struct {
	Name          string `json:"name"`
	qualifiedName string
	Namespaces    []*Namespace `json:"namespaces,omitempty"`
	Services      []*Service   `json:"services,omitempty"`
	DataTypes     []*DataType  `json:"dataTypes,omitempty"`
}

func (ns *Namespace) isRoot() bool {
	return ns.Name == ns.qualifiedName
}

func (ns *Namespace) addChild(child *Namespace) {
	ns.Namespaces = append(ns.Namespaces, child)
}

func (ns *Namespace) addService(service *Service) {
	ns.Services = append(ns.Services, service)
}

// DataTypeReference is a Name + Namespace
type DataTypeReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	// QualifiedName string `json:"qualifiedName"`
	// ElementType   *DataTypeReference `json:"elementType,omitempty"`
}

// RelativeDataTypeReference is a DataTypeReference with a relative path (used for referencing other files)
type RelativeDataTypeReference struct {
	DataTypeReference
	RelativePath string
	Alias        string
}

// DataType is a type
type DataType struct {
	DataTypeReference
	Base            *RelativeDataTypeReference   `json:"base,omitempty"`
	ReferencedTypes []*RelativeDataTypeReference `json:"referencedTypes,omitempty"`
	Fields          []*Field                     `json:"fields,omitempty"`
	Consts          []*Const                     `json:"consts,omitempty"`
	aliasMap        aliasMap
}

// Service is a type
type Service struct {
	DataType
	ServiceIdentifier    string
	Methods              []*Method `json:"methods"`
	ReferencedNamespaces []string
}

// Method within a Service
type Method struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeName string `json:"typeName"`
	Args     []*Arg `json:"args,omitempty"`
	nameSort string
}

// Arg is an argument to a Method within a Service
type Arg struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeName string `json:"typeName"`
}

// Field of a service
type Field struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	TypeName        string `json:"typeName"`
	ElementType     string `json:"elementType"`
	ElementTypeName string `json:"elementTypeName"`
}

// Const of a service
type Const struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	TypeName string      `json:"typeName"`
	Value    interface{} `json:"value"`
}

type aliasMap map[string]string

func (aliases aliasMap) nonEmpty() bool {
	return len(aliases) > 0
}

// ApplyAliases applies alias names to a method
func (meth *Method) ApplyAliases(aliases aliasMap) {
	if alias, contains := aliases[meth.Type]; contains {
		meth.TypeName = alias
	}

	for _, arg := range meth.Args {
		if alias, contains := aliases[arg.Type]; contains {
			arg.TypeName = alias
		}
	}
}

// ApplyAliases applies alias names to a field
func (field *Field) ApplyAliases(aliases aliasMap) {
	if alias, contains := aliases[field.Type]; contains {
		field.TypeName = alias
	}

	if alias, contains := aliases[field.ElementType]; contains {
		field.ElementTypeName = alias
	}
}

// ApplyAliases applies alias names to a const
func (con *Const) ApplyAliases(aliases aliasMap) {
	if alias, contains := aliases[con.Type]; contains {
		con.TypeName = alias
	}
}

// Aliaser is an interface for applying alias names to members
type Aliaser interface {
	ApplyAliases(aliasMap aliasMap)
}
