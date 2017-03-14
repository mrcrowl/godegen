package reflect

import "godegen/cli"

type Method struct {
	Name       string
	ReturnType *Type
	Signature  *MethodSig
}

func newMethod(methodDefRow *cli.MethodDefRow, asm *Assembly) *Method {
	sigBlob := methodDefRow.GetSignatureBlob()
	sigReader := NewSignatureReader(sigBlob, asm)
	paramRows := methodDefRow.GetParams(asm.metadata.Tables)

	return &Method{
		Name:       methodDefRow.Name,
		Signature:  sigReader.ReadMethodSignature(paramRows),
		ReturnType: nil,
	}
}
