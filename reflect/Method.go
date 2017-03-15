package reflect

import (
	"bytes"
	"godegen/cli"
)

type MemberAccess uint8

const (
	CompilerControlled MemberAccess = 0x00
	Private            MemberAccess = 0x01
	FamAndAssem        MemberAccess = 0x02
	Assem              MemberAccess = 0x03
	Family             MemberAccess = 0x04
	FamOrAssem         MemberAccess = 0x05
	Public             MemberAccess = 0x06
	MemberAccessMask   MemberAccess = 0x07
)

const (
	MethodAttributesStatic  uint16 = 0x10
	MethodAttributesFinal   uint16 = 0x20
	MethodAttributesVirtual uint16 = 0x40
)

var memberAccessNames = []string{"", "private", "", "internal", "protected", "", "public"}

type Method struct {
	name         string
	signature    *MethodSig
	memberAccess MemberAccess
	static       bool
	final        bool
	virtual      bool
}

func newMethod(methodDefRow *cli.MethodDefRow, asm *Assembly) *Method {
	sigBlob := methodDefRow.GetSignatureBlob()
	sigReader := NewSignatureReader(sigBlob, asm)
	paramRows := methodDefRow.GetParams(asm.metadata.Tables)
	memberAccess := MemberAccess(methodDefRow.Flags) & MemberAccessMask
	static := (methodDefRow.Flags & MethodAttributesStatic) > 0
	final := (methodDefRow.Flags & MethodAttributesFinal) > 0
	virtual := (methodDefRow.Flags & MethodAttributesVirtual) > 0

	return &Method{
		name:         methodDefRow.Name,
		signature:    sigReader.ReadMethodSignature(paramRows),
		memberAccess: memberAccess,
		static:       static,
		final:        final,
		virtual:      virtual,
	}
}

func (method *Method) Name() string {
	return method.name
}

func (method *Method) ReturnType() Type {
	return method.signature.returnType
}

func (method *Method) String() string {
	var buffer bytes.Buffer
	memberAccessName := memberAccessNames[method.memberAccess]
	if len(memberAccessName) > 0 {
		buffer.WriteString(memberAccessName)
		buffer.WriteByte(' ')
	}

	if method.static {
		buffer.WriteString("static ")
	}

	buffer.WriteString(method.ReturnType().FullName())
	buffer.WriteByte(' ')

	buffer.WriteString(method.Name())
	buffer.WriteByte('(')
	numParams := len(method.signature.parameters)
	if numParams > 0 {
		for i, param := range method.signature.parameters {
			buffer.WriteString(param.Type().FullName())
			buffer.WriteByte(' ')
			buffer.WriteString(param.Name())
			if (i + 1) < numParams {
				buffer.WriteString(", ")
			}
		}
	}
	buffer.WriteByte(')')
	return buffer.String()
}
