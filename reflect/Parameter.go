package reflect

type Parameter struct {
	name  string
	typ   Type
	flags uint16
}

func newParameter(name string, typ Type) *Parameter {
	return &Parameter{name, typ, 0}
}

func (param *Parameter) Name() string {
	return param.name
}

func (param *Parameter) Type() Type {
	return param.typ
}
