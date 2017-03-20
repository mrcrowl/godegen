package reflect

type Parameter struct {
	name  string
	typ   Type
	flags uint16
}

const (
	ParamAttributesIn         = 0x1
	ParamAttributesOut        = 0x2
	ParamAttributesOptional   = 0x10
	ParamAttributesHasDefault = 0x1000
)

func newParameter(name string, typ Type) *Parameter {
	return &Parameter{name, typ, 0}
}

func (param *Parameter) Name() string {
	return param.name
}

func (param *Parameter) Type() Type {
	return param.typ
}
