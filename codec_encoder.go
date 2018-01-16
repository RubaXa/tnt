package tnt

import (
	"bytes"

	"gopkg.in/vmihailenco/msgpack.v2"
	"gopkg.in/vmihailenco/msgpack.v2/codes"
)

type Encoder struct {
	writer  *bytes.Buffer
	encoder *msgpack.Encoder
	len     int
	arrLen  int
	err     error
}

func (e *Encoder) start() {
	e.writer.WriteByte(0)
}

func (e *Encoder) end() {
	l := e.len - e.arrLen

	if l < 16 {
		e.writer.Bytes()[0] = byte(codes.FixedArrayLow | byte(l))
	}
}

func (e *Encoder) IsDecode() bool {
	return false
}

func (e *Encoder) Int(v *int) {
	e.err = e.encoder.EncodeInt(*v)
	e.len++
}

func (e *Encoder) String(v *string) {
	e.err = e.encoder.EncodeString(*v)
	e.len++
}

func (e *Encoder) Slice(l int) int {
	e.err = e.encoder.EncodeSliceLen(l)
	e.arrLen += l
	return l
}
