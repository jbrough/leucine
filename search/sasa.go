package search

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/lib"
	"github.com/rs/zerolog/log"
)

// Parse FreeSASA residue results

func ProcessPdb(src, dst string) (err error) {
	paths, err := io.PathsFromOpt(src)
	if err != nil {
		return err
	}

	for _, path := range paths {
		dst_path := filepath.Join(dst, filepath.Base(path)+".sasa")
		if err = FreeSasa(path, dst_path); err != nil {
			return
		}

		log.Info().Msg("freesasa processed " + path)
	}

	return
}

func FreeSasa(src, dst string) (err error) {
	cmd := exec.Command(
		"freesasa", "--depth=residue", "--format=seq", "--output="+dst, src)

	if err = cmd.Run(); err != nil {
		return
	}

	return
}

func NewSasaSeqIndex() *SasaSeqIndex {
	return &SasaSeqIndex{
		seqs:   make(map[string]*SasaSeq),
		scores: make(map[string]int),
		Index:  NewIndex(),
	}
}

type SasaSeqIndex struct {
	seqs   map[string]*SasaSeq
	scores map[string]int
	Index  *Index
}

func (si *SasaSeqIndex) Get() (s *SasaSeq) {
	for _, s := range si.seqs {
		return s
	}
	return
}
func (si *SasaSeqIndex) Add(sa *SasaSeq) {
	si.seqs[sa.PdbId] = sa

	var seq string
	var score []int
	var lp int

	for i := 1; i < len(sa.seq); i++ {
		pos := sa.order[i]
		r := sa.seq[pos]
		if lp == 0 || pos-1 == lp {
			seq += r.AA.A
			score = append(score, int(r.Score))
		} else {
			seq += " "
			score = append(score, 0)
		}
		lp = pos
	}

	for i, word := range Words([]byte(seq), 20) {
		si.Index.AddVal([]byte(sa.PdbId), []byte(word), i)
	}
}

func NewSasaSeq(id string) *SasaSeq {
	return &SasaSeq{
		PdbId: id,
		seq:   make(map[int]*Residue),
		order: make(map[int]int),
	}
}

type SasaSeq struct {
	PdbId string
	seq   map[int]*Residue
	order map[int]int
	aa    string
}

func (s *SasaSeq) Region(a, b int) (string, string) {
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

func (se *SasaSeq) Add(r *Residue, i int) {
	se.seq[r.Pos-1] = r
	se.order[i] = r.Pos - 1
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

func LoadSasas(src string) (set *SasaSeqIndex, err error) {
	paths, err := io.PathsFromOpt(src)
	if err != nil {
		return
	}

	set = NewSasaSeqIndex()

	for _, path := range paths {
		sa, err := LoadSasa(path)
		if err != nil {
			return set, err
		}

		set.Add(sa)
	}

	return
}

func LoadSasa(path string) (sa *SasaSeq, err error) {
	aa := lib.AA()

	f, err := os.Open(path)
	if err != nil {
		return
	}

	id := strings.TrimSuffix(filepath.Base(path), ".ent.sasa")

	sa = NewSasaSeq(id)

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
			sa.Add(r, i-1)
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}
