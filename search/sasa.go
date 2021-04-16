package search

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/jbrough/leucine/lib"
)

// Parse FreeSASA residue results

func NewSasaSeqs() *SasaSeqs {
	return &SasaSeqs{
		make(map[int]*Residue),
		"",
	}
}

type SasaSeqs struct {
	seq map[int]*Residue
	aa  string
}

func (s *SasaSeqs) Region(a, b int) (string, string) {
	w := ""
	u := ""
	for i := a; i < b; i++ {
		r, ok := s.seq[i]
		if !ok {
			w += " "
			u += " "
		} else {
			n := int(r.Score / 10)
			if n > 9 {
				n = 9
			}
			v := strconv.Itoa(n)
			u += v
			if n > 5 {
				w += "x"
			} else {
				w += "y"
			}
		}
	}

	w = strings.Replace(w, "xx", "\u21F5", -1)
	w = strings.Replace(w, "y", "\u2509", -1)
	w = strings.Replace(w, "x", "\u2191", -1)

	return w, u
}

func (se *SasaSeqs) Add(r *Residue) {
	se.seq[r.Pos-1] = r
	se.aa += r.AA.A
}

type Residue struct {
	Seq   string
	Pos   int
	AA    *lib.AminoAcid
	Score float64
}

func parseResidue(b []byte, aa *lib.AminoAcids) (r *Residue, err error) {
	seq := string(b[4])
	pos, err := strconv.Atoi(strings.TrimSpace(string(b[6:11])))
	if err != nil {
		return
	}
	a, err := aa.GetAAA(string(b[12:15]))
	if err != nil {
		return nil, err
	}
	score, err := strconv.ParseFloat(strings.TrimSpace(string(b[20:25])), 64)
	if err != nil {
		return
	}
	r = &Residue{seq, pos, a, score}
	return
}

func LoadSasa(path string) (sa *SasaSeqs, err error) {
	aa := lib.AA()

	f, err := os.Open(path)
	if err != nil {
		return
	}

	sa = NewSasaSeqs()

	scanner := bufio.NewScanner(f)
	var i int
	for scanner.Scan() {
		i++
		l := scanner.Bytes()
		if i == 1 {
			continue
		}
		if len(l) == 0 {
			continue
		}
		var r *Residue
		tmp := make([]byte, len(l))
		copy(tmp, l)
		r, err = parseResidue(tmp, aa)
		if err != nil {
			return
		}
		if r.Seq == "A" {
			sa.Add(r)
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}
