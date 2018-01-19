package tnt

import "testing"

type MySpace struct {
	Space
}

func (s *MySpace) NextEntry() IEntry {
	return &Record{}
}

func (s *MySpace) Select(c Cursor) ([]*Record, error) {
	list := make([]*Record, 0)
	err := s.exec(c).Each(func(e IEntry, _ int) {
		list = append(list, e.(*Record))
	})
	return list, err
}

func (s *MySpace) SelectOne(c Cursor) (*Record, error) {
	c.Iterator = IterEq
	entry, err := s.exec(c).First()
	return entry.(*Record), err
}

var spaceFactory = func() (Index, ISpace) {
	return Index("primary"), &MySpace{}
}

func (b *MyBox) GetTestSpace() *MySpace {
	return b.InitSpace("tnt_space", spaceFactory).(*MySpace)
}

func TestSpaceGetAndSelect(t *testing.T) {
	box, _ := GetBox("test", myBoxFactory).(*MyBox)
	space := box.GetTestSpace()

	if space.Name != "tnt_space" {
		t.Errorf("space.Name: %s != tnt_space", space.Name)
	}

	if space.DefaultIdx != "primary" {
		t.Errorf("space.DefaultIdx: %s != primary", space.DefaultIdx)
	}

	records, err := space.Select(Cursor{})

	if err != nil {
		t.Error("space.Select().err:", err)
	}

	if len(records) != 2 {
		t.Errorf("len(records): %d != 2", len(records))
	}
}

func TestSpaceSelecOne(t *testing.T) {
	box, _ := GetBox("test", myBoxFactory).(*MyBox)
	space := box.GetTestSpace()

	rec, err := space.SelectOne(Cursor{
		Key: 2,
	})

	if err != nil {
		t.Error("space.SelectOne().err:", err)
	}

	if !rec.IsExists() {
		t.Error("space.SelectOne(): not exists")
	}

	if rec.Id != 2 {
		t.Errorf("space.SelectOne().Id: %d != 2", rec.Id)
	}
}

func TestSpaceSelecOneNotExists(t *testing.T) {
	box, _ := GetBox("test", myBoxFactory).(*MyBox)
	space := box.GetTestSpace()

	rec, err := space.SelectOne(Cursor{
		Key: 666,
	})

	if err != nil {
		t.Error("space.SelectOne().err:", err)
	}

	if rec.IsExists() {
		t.Error("space.SelectOne(): exists?!")
	}

	if rec.Id != 0 {
		t.Errorf("space.SelectOne().Id: %d != 0", rec.Id)
	}
}
