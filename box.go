package tnt

import (
	"fmt"
	"math"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

type Config struct {
	Server  string
	User    string
	Pass    string
	Timeout time.Duration
}

type IBox interface {
	init(ns string, cfg *Config)
	Select(s string, c Cursor, next NextEntry) *LazyResult
}

type Box struct {
	IBox
	ns     string
	cfg    *Config
	conn   *tarantool.Connection
	err    error
	spaces *registry
}

func (b *Box) init(ns string, cfg *Config) {
	b.ns = ns
	b.cfg = cfg
	b.spaces = &registry{}

	if cfg.Timeout == 0 {
		cfg.Timeout = 500
	}

	opts := tarantool.Opts{
		Timeout:       cfg.Timeout * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          cfg.User,
		Pass:          cfg.Pass,
	}

	b.conn, b.err = tarantool.Connect(cfg.Server, opts)
}

func (b *Box) InitSpace(n string, f func() (Index, ISpace)) ISpace {
	return b.spaces.Get(n, func() interface{} {
		idx, space := f()
		space.init(n, idx, b, space.NextEntry)
		return space
	}).(ISpace)
}

func (b *Box) Insert(s string, v interface{}) {
	res, err := b.conn.Insert(s, v)
	fmt.Printf("Err: %#v\n", err)
	fmt.Printf("Res: %#v\n", res)
}

func (b *Box) Select(s string, c Cursor, next NextEntry) *LazyResult {
	if b.err != nil {
		return &LazyResult{err: b.err}
	}

	if c.All {
		if c.Limit == 0 {
			c.Limit = math.MaxUint32
		}

		if c.Iterator == 0 {
			c.Iterator = IterAll
		}
	}

	if c.Iterator == IterEq {
		c.Offset = 0
		c.Limit = 1
	}

	where := []interface{}{}

	if c.Where != nil {
		where = append(where, c.Where)
	}

	future := b.conn.SelectAsync(
		s,
		string(c.Index),
		c.Offset,
		c.Limit,
		c.Iterator,
		where,
	)

	return &LazyResult{
		future: future,
		next:   next,
	}
}

var boxes = &registry{}

func GetBox(ns string, f func() (IBox, *Config)) IBox {
	return boxes.Get(ns, func() interface{} {
		box, cfg := f()
		box.init(ns, cfg)
		return box
	}).(IBox)
}
