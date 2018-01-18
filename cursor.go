package tnt

type Cursor struct {
	Index    Index
	Iterator uint32
	Where    interface{}
	All      bool
	Offset   uint32
	Limit    uint32
}
