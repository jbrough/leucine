package search

import (
	"bufio"
	"errors"

	"github.com/twmb/murmur3"
)

func NewIndex() *Index {
	return &Index{
		Test:   make(map[uint64]bool),
		Match:  make(map[uint64]map[uint64][]int),
		GetKey: make(map[uint64]string),
		GetRef: make(map[uint64][2][]byte),
	}
}

type Index struct {
	Test   map[uint64]bool
	Match  map[uint64]map[uint64][]int
	GetKey map[uint64]string
	GetRef map[uint64][2][]byte
}

func (i *Index) Hash(data []byte) uint64 {
	hasher := murmur3.New128()
	hasher.Write(data)
	v1, v2 := hasher.Sum128()
	_ = v2
	return v1
}

func (i *Index) AddKey(b []byte) bool {
	h := i.Hash(b)
	// return false if key present to avoid collisions
	if _, ok := i.GetKey[h]; ok {
		return false
	}

	i.GetKey[h] = string(b)
	i.Match[h] = make(map[uint64][]int)

	return true
}

func (i *Index) AddVal(key, val []byte, idx int) bool {
	kh := i.Hash(key)
	vh := i.Hash(val)

	// return false if key not present
	if _, ok := i.GetKey[kh]; !ok {
		return false
	}

	// return false if collision
	if _, ok := i.Test[vh]; ok {
		return false
	}

	i.Test[vh] = true
	idxs := i.Match[kh][vh]
	i.Match[kh][vh] = append(idxs, idx)

	return true
}

func (i *Index) AddRef(key []byte, val [2][]byte) {
	h := i.Hash(key)

	i.GetRef[h] = val
}

func IndexStream(scanner *bufio.Scanner, ngram_n int) (index *Index, err error) {
	index = NewIndex()

	d := true
	var def []byte
	var seq []byte
	for scanner.Scan() {
		l := scanner.Bytes()
		if d {
			d = !d

			def = nil
			def = make([]byte, len(l)-1) // remove leading '>'
			copy(def, l[1:])

		} else {
			if ok := index.AddKey(def); !ok {
				err = errors.New("key exists: " + string(def))
				return
			}

			for i, word := range words(l, ngram_n) {
				index.AddVal(def, word, i)

				seq = nil
				seq = make([]byte, len(l))
				copy(seq, l)

				index.AddRef(def, [2][]byte{def, seq})
			}

			d = !d
		}
	}

	err = scanner.Err()

	return
}
