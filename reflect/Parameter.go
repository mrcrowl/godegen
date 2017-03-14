package reflect

type Parameter struct {
	name string
	typ  Type
}

func newParameter(name string, typ Type) *Parameter {
	return &Parameter{name, typ}
}

func (param *Parameter) Name() string {
	return param.name
}

func (param *Parameter) Type() Type {
	return param.typ
}
