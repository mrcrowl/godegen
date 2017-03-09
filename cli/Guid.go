package cli

type Guid [16]byte

func ZeroGuid() Guid {
	return Guid{}
}

func NewGuid(bytes []byte) Guid {
	if len(bytes) != 16 {
		panic("Guid must be 16 bytes")
	}

	guid := Guid{}
	copy(guid[:], bytes)
	return guid
}
