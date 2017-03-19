package description

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
	QualifiedName string       `json:"qualifiedName"`
	Namespaces    []*Namespace `json:"namespaces,omitempty"`
	Services      []*Service   `json:"services,omitempty"`
	DataTypes     []*DataType  `json:"dataTypes,omitempty"`
}

func (ns *Namespace) isRoot() bool {
	return ns.Name == ns.QualifiedName
}

func (ns *Namespace) addChild(child *Namespace) {
	ns.Namespaces = append(ns.Namespaces, child)
}

func (ns *Namespace) addService(service *Service) {
	ns.Services = append(ns.Services, service)
}

type DataTypeReference struct {
	Name          string             `json:"name"`
	Namespace     string             `json:"namespace"`
	QualifiedName string             `json:"qualifiedName"`
	ElementType   *DataTypeReference `json:"elementType,omitempty"`
}

type Service struct {
	DataTypeReference
	Methods []*Method `json:"methods"`
}

type DataType struct {
	DataTypeReference
	Base   *DataTypeReference `json:"base,omitempty"`
	Fields []*Field           `json:"fields,omitempty"`
}

type Method struct {
	Name string             `json:"name"`
	Type *DataTypeReference `json:"type"`
	Args []*Arg             `json:"args,omitempty"`
}

type Arg struct {
	Name string             `json:"name"`
	Type *DataTypeReference `json:"type"`
}

type Field struct {
	Name string             `json:"name"`
	Type *DataTypeReference `json:"type"`
}

type Const struct {
	Name  string             `json:"name"`
	Type  *DataTypeReference `json:"type"`
	Value interface{}        `json:"value"`
}
