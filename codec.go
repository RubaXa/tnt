package tnt

import (
	"bytes"

	"gopkg.in/vmihailenco/msgpack.v2"
)

type Codec interface {
	IsDecode() bool

	Int(v *int)
	String(v *string)
	Slice(l int) int
}

type ICodecSupported interface {
	Codec(e Codec)
}

func Encode(v ICodecSupported) ([]byte, error) {
	w := &bytes.Buffer{}
	e := &Encoder{
		writer:  w,
		encoder: msgpack.NewEncoder(w),
	}

	e.start()
	v.Codec(e)
	e.end()

	if e.err != nil {
		return nil, e.err
	}

	return w.Bytes(), nil
}

func Decode(t ICodecSupported, s []byte) error {
	d := &Decoder{
		decoder: msgpack.NewDecoder(bytes.NewReader(s)),
	}

	d.start()
	t.Codec(d)
	d.end()

	return d.err
}

func DecodeEntry(md *msgpack.Decoder, e IEntry) error {
	d := &Decoder{
		decoder: md,
	}

	d.start()
	e.Codec(d)
	d.end()

	return d.err
}
