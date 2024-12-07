package hash

import (
	"encoding/binary"
	"fmt"
	"github.com/dchest/siphash"
	"math"
)

const (
	DictHtInitialSize = 4
	DictOK            = 0
	DictErr           = 1
)

/* ---------- Type Definition ----------*/
type dictEntry struct {
	key   string
	value interface{}
	next  *dictEntry
}

type dictHt struct {
	table    []*dictEntry
	size     uint
	sizemask uint
	used     uint
}

type dictIterator struct {
	dict      *Dict
	index     int
	entry     *dictEntry
	nextEntry *dictEntry
}

type Dict struct {
	ht        [2]dictHt
	rehashIdx int
	iterator  int
}

/* ---------- Hash Function ----------*/
var dictHashFunctionSeed [16]byte

func dictSetHashFunctionSeed(seed [16]byte) {
	dictHashFunctionSeed = seed
}

func dictGetHashFunction() [16]byte {
	return dictHashFunctionSeed
}

/* ---------- API Implementation ----------*/

func initDict() *Dict {
	return &Dict{
		ht: [2]dictHt{
			{table: nil, size: 0, sizemask: 0, used: 0},
			{table: nil, size: 0, sizemask: 0, used: 0},
		},
		rehashIdx: -1,
		iterator:  0,
	}
}

func dictResize(dict *Dict) int {
	if dictIsRehashing(dict) {
		return DictErr
	}
	minimal := dict.ht[0].used
	if minimal < DictHtInitialSize {
		minimal = DictHtInitialSize
	}

	var newHt dictHt
	realsize, err := dictNextPower(uint(minimal))
	if err != nil {
		return DictErr
	}
	newHt.size = realsize
	newHt.sizemask = realsize - 1
	newHt.table = make([]*dictEntry, newHt.size)
	newHt.used = 0

	//first initialization
	if dict.ht[0].table == nil {
		dict.ht[0] = newHt
		return DictOK
	}

	dict.ht[1] = newHt
	dict.rehashIdx = 0
	return DictOK
}

func dictIsRehashing(dict *Dict) bool {
	return dict.rehashIdx != -1
}

func dictNextPower(size uint) (uint, error) {
	if size > math.MaxUint32 {
		return 0, fmt.Errorf("dict size is too big: %d", size)
	}
	i := DictHtInitialSize
	for {
		if uint(i) >= size/2 {
			return uint(i), nil
		}
		if uint(i) >= size {
			return uint(i), nil
		}
		i *= 2
	}
}

func dictRehash(dict *Dict, n int) int {
	emptyVisits := n * 10
	if !dictIsRehashing(dict) {
		return 0
	}

	for i := 0; i < n && dict.ht[0].used > 0; i++ {
		if dict.ht[0].size < uint(dict.rehashIdx) {
			return DictErr
		}
		//skip empty bucket
		for dict.ht[0].table[dict.rehashIdx] == nil {
			dict.rehashIdx++
			emptyVisits--
			if emptyVisits <= 0 {
				return DictErr
			}
		}

		de := dict.ht[0].table[dict.rehashIdx]
		for de != nil {
			nextde := de.next
			seed := dictGetHashFunction()
			k0 := binary.LittleEndian.Uint64(seed[0:8])
			k1 := binary.LittleEndian.Uint64(seed[8:16])
			hashValue := siphash.Hash(k0, k1, []byte(de.key))
			de.next = dict.ht[1].table[hashValue]
			dict.ht[1].table[hashValue] = de
			dict.ht[0].used--
			dict.ht[1].used++
			de = nextde
		}
		dict.ht[0].table[dict.rehashIdx] = nil
		dict.rehashIdx++
	}

	if dict.ht[0].used == 0 {
		dict.ht[0].table = nil
		dict.ht[0] = dict.ht[1]
		dict.ht[1] = dictHt{
			table:    nil,
			size:     0,
			sizemask: 0,
			used:     0,
		}
		dict.rehashIdx = -1
		return DictOK
	}
	return DictErr
}
