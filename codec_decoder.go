package tnt

import (
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Decoder struct {
	decoder *msgpack.Decoder
	err     error
}

func (d *Decoder) start() {
	_, d.err = d.decoder.DecodeSliceLen()
}

func (d *Decoder) end() {
}

func (d *Decoder) IsDecode() bool {
	return true
}

func (d *Decoder) Int(v *int) {
	*v, d.err = d.decoder.DecodeInt()
}

func (d *Decoder) String(v *string) {
	*v, d.err = d.decoder.DecodeString()
}

func (d *Decoder) Slice(l int) int {
	l, d.err = d.decoder.DecodeSliceLen()
	return l
}
