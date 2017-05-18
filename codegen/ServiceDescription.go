package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ServiceDescription struct {
	Namespaces []*Namespace
}

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

type Namespace struct {
	Name          string       `json:"name"`
	qualifiedName string       `json:"-"`
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

type DataTypeReference struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	// QualifiedName string `json:"qualifiedName"`
	// ElementType   *DataTypeReference `json:"elementType,omitempty"`
}

type MappedDataTypeReference struct {
	DataTypeReference
	MappedType string
}

type DataType struct {
	DataTypeReference
	Base   *DataTypeReference `json:"base,omitempty"`
	Fields []*Field           `json:"fields,omitempty"`
	Consts []*Const           `json:"consts,omitempty"`
}

type Service struct {
	DataType
	ServiceIdentifier    string
	Methods              []*Method `json:"methods"`
	ReferencedNamespaces []string
}

type Method struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeName string `json:"typeName"`
	Args     []*Arg `json:"args,omitempty"`
	nameSort string
}

type Arg struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeName string `json:"typeName"`
}

type Field struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	TypeName        string `json:"typeName"`
	ElementType     string `json:"elementType"`
	ElementTypeName string `json:"elementTypeName"`
}

type Const struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	TypeName string      `json:"typeName"`
	Value    interface{} `json:"value"`
}
