package tnt

const (
	f_ENTRY_EXISTS  = 1 << 1
	f_DECODE_FAILED = 1 << 2
)

type NextEntry func() IEntry
type EntryIterator func(e IEntry, i int)

type IEntry interface {
	ICodecSupported
	IsExists() bool
	addEntryFlag(f byte)
}

type Entry struct {
	entryFlags byte
}

func (e *Entry) addEntryFlag(f byte) {
	e.entryFlags |= f
}

func (e *Entry) IsExists() bool {
	return (e.entryFlags & f_ENTRY_EXISTS) != 0
}
