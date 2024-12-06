package hash

const (
	DICT_HT_INITIAL_SIZE = 4
)

type dictEntry struct {
	key   string
	value interface{}
	next  *dictEntry
}

type dictHt struct {
	table    []*dictEntry
	size     int
	sizemask int
	used     int
}

type dictIterator struct {
	dict      *dict
	index     int
	entry     *dictEntry
	nextEntry *dictEntry
}

type dict struct {
	ht        [2]dictHt
	rehashIdx int
	iterator  *dictIterator
}
