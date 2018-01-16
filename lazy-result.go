package tnt

import (
	tarantool "github.com/tarantool/go-tarantool"
	msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

type LazyResult struct {
	future    *tarantool.Future
	next      NextEntry
	interator EntryIterator
	entries   []IEntry
	parsed    bool
	err       error
	size      int
}

func (lr *LazyResult) DecodeMsgpack(d *msgpack.Decoder) (err error) {
	lr.size, err = d.DecodeSliceLen()

	if err == nil {
		for i := 0; i < lr.size; i++ {
			entry := lr.next()
			err = DecodeEntry(d, entry)
			lr.entries = append(lr.entries, entry)

			if err != nil {
				entry.addEntryFlag(f_DECODE_FAILED)
			} else {
				entry.addEntryFlag(f_ENTRY_EXISTS)
			}

			if lr.interator != nil {
				lr.interator(entry, i)
			}
		}
	}

	return
}

func (lr *LazyResult) parse() error {
	if !lr.parsed {
		lr.parsed = true
		lr.err = lr.future.GetTyped(lr)
	}

	return lr.err
}

func (lr *LazyResult) Error() error {
	return lr.parse()
}

func (lr *LazyResult) Entries() ([]IEntry, error) {
	err := lr.parse()
	return lr.entries, err
}

func (lr *LazyResult) Size() (int, error) {
	err := lr.parse()
	return lr.size, err
}

func (lr *LazyResult) Each(fn EntryIterator) error {
	if lr.parsed {
		for i, e := range lr.entries {
			fn(e, i)
		}
	} else {
		lr.interator = fn
		lr.parse()
		lr.interator = nil
	}

	return lr.err
}

func (lr *LazyResult) First() (IEntry, error) {
	err := lr.parse()

	if lr.size > 0 {
		return lr.entries[0], nil
	} else {
		return lr.next(), err
	}
}

func (lr *LazyResult) Last() (IEntry, error) {
	err := lr.parse()

	if lr.size > 0 {
		return lr.entries[lr.size-1], nil
	} else {
		return lr.next(), err
	}
}
