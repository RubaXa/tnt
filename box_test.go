package tnt

import (
	"testing"
)

type MyBox struct {
	Box
}

type Record struct {
	Entry
	Id   int
	Name string
}

func (r *Record) Codec(c Codec) {
	c.Int(&r.Id)
	c.String(&r.Name)
}

var myBoxFactory = func() (IBox, *Config) {
	return &MyBox{}, &Config{
		Server: "localhost:1311",
		User:   "tnt_test",
		Pass:   "tnt_test",
	}
}

func TestGetBox(t *testing.T) {
	box, ok := GetBox("test", myBoxFactory).(*MyBox)

	if !ok {
		t.Error("MyBox failed")
	}

	if box.ns != "test" {
		t.Errorf("Wrong ns: %s", box.ns)
	}

	if box.cfg.Server != "localhost:1311" {
		t.Errorf("Wrong cfg: %#v", box.cfg)
	}

	if box.err != nil {
		t.Errorf("Box.err: %s, %#v", box.err, box.cfg)
	}
}

func TestBoxSelect(t *testing.T) {
	box := GetBox("test", myBoxFactory).(*MyBox)
	next := func() IEntry { return &Record{} }
	res := box.Select("tnt_space", Cursor{
		Index: "primary",
		All:   true,
	}, next)

	list := make([]*Record, 0)
	err := res.Each(func(entry IEntry, _ int) {
		list = append(list, entry.(*Record))
	})

	if err != nil {
		t.Error("Each.err:", err)
	}

	if len(list) != 2 {
		t.Errorf("Each.len: %d != 2", len(list))
	}

	if list[0] == list[1] {
		t.Errorf("Entries must be not equal: %#v and %v", list[0], list[1])
	}

	entry, err := res.First()
	if err != nil {
		t.Error("First.err:", err)
	}

	if entry == nil {
		t.Fatal("first 'entry' is nil")
	}

	if !entry.IsExists() {
		t.Error("first 'entry' not exists")
	}

	first := entry.(*Record)
	if first.Id != 1 {
		t.Errorf("first.Id: %d != 1", first.Id)
	}

	entry, _ = res.Last()
	last := entry.(*Record)
	if last.Id != 2 {
		t.Errorf("last.Id: %d != 2", last.Id)
	}

	if first == last {
		t.Error("first == last")
	}
}
