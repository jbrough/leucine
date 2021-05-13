package seq

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/jbrough/leucine/lib"
)

func NewDict() *Dict {
	idm := make(map[string]uint16)
	tgm := make(map[uint16]string)
	aa := strings.ToUpper(lib.ProteinAlphabet)

	var i uint16
	for _, co := range aa {
		idm[string(co)] = i
		tgm[i] = string(co)
		i++
	}

	for n := 2; n < 4; n++ {
		for w := range words(aa, n) {
			if len(w) == n {
				idm[w] = i
				tgm[i] = w
				i++
			}
		}
	}

	return &Dict{
		idm, tgm,
	}
}

type Dict struct {
	id map[string]uint16
	tg map[uint16]string
}

func (d Dict) Id(tr string) uint16 {
	return d.id[tr]
}

func (d Dict) Tg(id uint16) string {
	return d.tg[id]
}

func Bitmap(seq string, dict *Dict) *roaring.Bitmap {
	n := 3

	bm := roaring.New()
	for i, _ := range seq {
		if i+n > len(seq) {
			break
		} else {
			k := seq[i : i+n]
			if v, ok := dict.id[k]; ok {
				bm.Add(uint32(v))
			}
		}
	}

	return bm
}

func Pack(s string, dict *Dict, buf *bytes.Buffer) []byte {
	var keys []uint16

	for _, w := range chunks(s, 3) {
		keys = append(keys, dict.Id(w))
	}

	if err := binary.Write(buf, binary.LittleEndian, keys); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func Unpack(buf *bytes.Buffer, dict *Dict) string {
	keys := make([]uint16, buf.Len()/2)
	if err := binary.Read(buf, binary.LittleEndian, &keys); err != nil {
		panic(err)
	}

	var s string
	for i := 0; i < len(keys); i++ {
		s += dict.Tg(keys[i])
	}

	return s
}

func chunks(seq string, n int) (r []string) {
	x := 0
	for _ = range seq {
		if x+n >= len(seq) {
			r = append(r, seq[x:])
			break
		} else {
			r = append(r, seq[x:x+n])
		}

		x += n
	}

	return
}

func words(alphabet string, length int) <-chan string {
	c := make(chan string)

	go func(c chan string) {
		defer close(c)
		addLetter(c, "", alphabet, length)
	}(c)

	return c
}

func addLetter(c chan string, combo string, alphabet string, length int) {
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}
