package reflect

import "godegen/cli"

type Method struct {
	Name       string
	ReturnType *Type
	Blob       cli.Blob
}

func newMethod(methodDefRow *cli.MethodDefRow, asm *Assembly) *Method {
	return &Method{
		Name:       methodDefRow.Name,
		Blob:       methodDefRow.GetSignature(),
		ReturnType: nil,
	}
}
