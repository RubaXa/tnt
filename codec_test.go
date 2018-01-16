package tnt

import (
	"testing"

	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type Mock struct {
	Id    int
	Name  string
	Flags []string
}

func (m *Mock) EncodeMsgpack(e *msgpack.Encoder) error {
	e.EncodeSliceLen(2)
	e.EncodeInt(m.Id)
	e.EncodeString(m.Name)

	e.EncodeSliceLen(len(m.Flags))
	for _, v := range m.Flags {
		e.EncodeString(v)
	}

	return nil
}

func (m *Mock) Codec(c Codec) {
	c.Int(&m.Id)
	c.String(&m.Name)

	if l := c.Slice(len(m.Flags)); c.IsDecode() {
		m.Flags = make([]string, l)
	}

	for i := range m.Flags {
		c.String(&m.Flags[i])
	}
}

func TestEncode(t *testing.T) {
	mock := &Mock{123, "Test", []string{"Foo", "Bar"}}
	actual, err := Encode(mock)

	if err != nil {
		t.Errorf("Error: %#v", err)
		return
	}

	if actual == nil {
		t.Error("actual must be not nil")
		return
	}

	expected, err := msgpack.Marshal(mock)
	if err != nil {
		panic(err)
	}

	al := len(actual)
	el := len(expected)

	if al != el {
		t.Errorf("Length: %d != %d", al, el)
		return
	}

	for i, a := range actual {
		if a != expected[i] {
			t.Errorf("Byte #%d: %d != %d", i, a, expected[i])
		}
	}
}

func TestDecode(t *testing.T) {
	expected := &Mock{123, "Test", []string{"Foo", "Bar"}}
	expectedBytes, err := msgpack.Marshal(expected)

	if err != nil {
		t.Errorf("expected (encode err): %#v", err)
		return
	}

	var actual Mock
	err = Decode(&actual, expectedBytes)

	if err != nil {
		t.Errorf("Decode err: %#v", err)
		return
	}

	if actual.Id != expected.Id {
		t.Errorf("actual.Id: %d != %d", actual.Id, expected.Id)
	}

	if actual.Name != expected.Name {
		t.Errorf("actual.Name: %s != %s", actual.Name, expected.Name)
	}

	if len(actual.Flags) != len(expected.Flags) {
		t.Errorf("actual.Flags.len: %d != %d", len(actual.Flags), len(expected.Flags))
	}

	for i, f := range actual.Flags {
		if f != expected.Flags[i] {
			t.Errorf("actual.Flags[%d]: %v != %v", i, f, expected.Flags[i])
		}
	}
}
