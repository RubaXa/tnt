package tnt

import (
	"fmt"
	"testing"
)

const PREFIX = "item:%d"

type Item struct {
	cid int
}

func (i Item) String() string {
	return fmt.Sprintf(PREFIX, i.cid)
}

func TestGet(t *testing.T) {
	reg := &registry{}

	expected, ok := reg.Get("item", func() interface{} {
		return Item{1}
	}).(Item)

	if !ok {
		t.Error("Casting failed")
	}

	actual, ok := reg.Get("item", func() interface{} {
		return Item{2}
	}).(Item)

	if expected != actual {
		t.Error("Expected != Actual")
	}

	if actual.cid != 1 {
		t.Error("Failed actual.cid")
	}
}
