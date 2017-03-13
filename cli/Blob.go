package cli

type Blob struct {
	Length uint32
	Data   []byte
}

func NewBlob(length uint32, buffer []byte) *Blob {
	return &Blob{length, buffer}
}

func ZeroBlob() *Blob {
	return &Blob{}
}
