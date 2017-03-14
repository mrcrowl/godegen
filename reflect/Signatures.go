package reflect

import "godegen/cli"
import "bytes"

const (
	DEFAULT      = 0x0
	VARARG       = 0x5
	GENERIC      = 0x10
	HASTHIS      = 0x20
	EXPLICITTHIS = 0x40
)

func ReadMethodSignature(blob cli.Blob, paramRows []*cli.ParamRow) *MethodSig {
	reader := cli.NewShapeReader(bytes.NewReader(blob.Data))

	flags := reader.ReadByte()
	// hasExplicitThis := flags&(HASTHIS|EXPLICITTHIS) > 0
	// vararg := (flags & VARARG) > 0
	generic := (flags & GENERIC) > 0
	// genParamCount := 0
	if generic {
		// genParamCount := reader.ReadCompressedUInt()
	}

	// paramCount := reader.ReadCompressedUInt()
	return nil
}

type MethodSig struct {
	returnType Type
	parameters []Parameter
}
