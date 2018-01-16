package tnt

type Cursor struct {
	Index    Index
	Iterator uint32
	All      bool
	Offset   uint32
	Limit    uint32
}
