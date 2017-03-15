package cli

type IRow interface {
	RowNumber() uint32
	String() string
}
