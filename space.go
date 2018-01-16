package tnt

import "fmt"

type Index string

type ISpace interface {
	init(name string, idx Index, box IBox, nentry NextEntry)
	NextEntry() IEntry
}

type Space struct {
	ISpace
	Name       string
	DefaultIdx Index
	box        IBox
	nextEntry  NextEntry
}

func (s *Space) init(name string, idx Index, box IBox, nentry NextEntry) {
	s.box = box
	s.Name = name
	s.DefaultIdx = idx
	s.nextEntry = nentry
}

func (s *Space) exec(c Cursor) *LazyResult {
	if c.Iterator == 0 && c.Limit == 0 {
		c.All = true
	}

	if s.box == nil {
		return &LazyResult{
			err: fmt.Errorf("[tnt] Space.'%s' not inited", s.Name),
		}
	}

	if s.nextEntry == nil {
		return &LazyResult{
			err: fmt.Errorf("[tnt] Space.'%s'.nextEntry is nil", s.Name),
		}
	}

	if c.Index == "" {
		c.Index = s.DefaultIdx
	}

	if c.Iterator == 0 {
		c.Iterator = IterAll
	}

	return s.box.Select(s.Name, c, s.nextEntry)
}
