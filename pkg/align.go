package blastr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/twmb/murmur3"
)

type AlignInfo struct {
	Query      string     `json:"query"`
	Candidates string     `json:"candidates";`
	Stats      AlignStats `json:"stats"`
}

type AlignStats struct {
	SequencesSearched    int          `json:"sequeneces_searched"`
	AlignmentsFound      int          `json:"alignments_found"`
	AlignmentsTested     uint64       `json:"alignments_tested"`
	AlignmentTestsPerSec uint64       `json:"alginment_tests_per_sec"`
	RuntimeSecs          float64      `json:"runtime_secs"`
	FastaFile            string       `json:"fasta_file,omitempty"`
	Stats                []AlignStats `json:"stats,omitempty"`
}

func (s *AlignStats) Add(sa AlignStats) {
	s.SequencesSearched += sa.SequencesSearched
	s.AlignmentsFound += sa.AlignmentsFound
	s.AlignmentsTested += sa.AlignmentsTested
	s.AlignmentTestsPerSec += sa.AlignmentTestsPerSec
	s.RuntimeSecs += sa.RuntimeSecs
}

func (s AlignStats) AsJSON() string {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(j)
}

type Alignment struct {
	QueryId     string         `json:"qid"`
	QueryIdx    int            `json:"qi"`
	SubjectId   string         `json:"sid"`
	SubjectIdx  int            `json:"si"`
	Word        string         `json:"w"`
	QuerySeq    LocalAlignment `json:"qs,omitempty"`
	SubjectSeq  LocalAlignment `json:"ss,omitempty"`
	QueryName   string         `json:"qn,omitempty"`
	SubjectName string         `json:"sn,omitempty"`
}

type LocalAlignment struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	A string `json:"a"`
}

func localAligments(ba, bb []byte, ia, ib, sa, sb int) (LocalAlignment, LocalAlignment) {
	seq := func(b []byte, i, s int) ([]byte, int, int) {
		if s%2 != 0 {
			s++
		}
		n := (60 - s) / 2

		x := i - n
		if x < 0 {
			x = 0
		}

		y := i + s + n

		if y >= len(b) {
			y = len(b)
		}

		return b[x:y], x, y
	}

	aseq, ax, ay := seq(ba, ia, sa)
	bseq, bx, by := seq(bb, ib, sb)

	var as, bs string
	if ay >= len(ba)-1 {
		as += "*"
	}

	if by >= len(bb)-1 {
		bs += "*"
	}

	if len(aseq) < len(bseq) {
		bseq = bseq[:len(aseq)]
	} else if len(bseq) < len(aseq) {
		aseq = aseq[:len(bseq)]
	}

	a := string(aseq) + as
	b := string(bseq) + bs

	return LocalAlignment{
			X: ax,
			Y: ay,
			A: a,
		}, LocalAlignment{
			X: bx,
			Y: by,
			A: b,
		}
}

func hash(data []byte) uint64 {
	hasher := murmur3.New128()
	hasher.Write(data)
	v1, v2 := hasher.Sum128()
	_ = v2
	return v1
}

func words(seq []byte, n int) (r [][]byte) {
	s := len(seq)

	for i, _ := range seq {
		if i+n > s {
			return append(r, seq[i:])
		} else {
			r = append(r, seq[i:i+n])
		}
	}

	return r
}

func btoid(a []byte) []byte {
	i := bytes.IndexByte(a[4:], '|')
	return a[4 : i+4]
}

func Align(query_path, test_path string, ngram_n int, out chan Alignment) (stats AlignStats, err error) {
	query_file, err := os.Open(query_path)
	if err != nil {
		return
	}

	query_test := make(map[uint64]bool)
	query_table := make(map[uint64]map[uint64]int)
	query_ids := make(map[uint64]string)
	query_detail := make(map[uint64][2][]byte)
	stats = AlignStats{}

	d := true
	var id []byte
	var b []byte
	scanner := bufio.NewScanner(query_file)
	for scanner.Scan() {
		l := scanner.Bytes()
		if d {
			id = btoid(l)
			b = l
			d = !d
		} else {
			query_index := make(map[uint64]int)
			for i, word := range words(l, ngram_n) {
				h := hash(word)
				query_test[h] = true
				query_index[h] = i
			}
			query_table[hash(id)] = query_index
			query_ids[hash(id)] = string(id)
			query_detail[hash(id)] = [2][]byte{b, l}
			d = !d
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	query_file.Close()

	test_file, err := os.Open(test_path)
	if err != nil {
		return
	}
	defer test_file.Close()

	scanner = bufio.NewScanner(test_file)

	ts := time.Now()

	for scanner.Scan() {
		l := scanner.Bytes()
		var skip int
		if d {
			stats.SequencesSearched++
			id = btoid(l)
			b = l
			d = !d
		} else {
			// dont check self
			if _, ok := query_table[hash(id)]; !ok {
				for i, word := range words(l, ngram_n) {
					if skip > 0 && i < skip {
						continue
					} else {
						skip = 0
					}
					stats.AlignmentsTested++
					if _, ok := query_test[hash(word)]; ok {
						for qid, tbl := range query_table {
							if idx, ok := tbl[hash(word)]; ok {
								skip = i + ngram_n

								sseq, qseq := localAligments(l, query_detail[qid][1], i, idx, ngram_n, ngram_n)

								out <- Alignment{
									QueryId:     query_ids[qid],
									QueryIdx:    idx,
									SubjectId:   string(id),
									SubjectIdx:  i,
									Word:        string(word),
									QuerySeq:    qseq,
									SubjectSeq:  sseq,
									QueryName:   string(query_detail[qid][0][1:]),
									SubjectName: string(b[1:]),
								}

								stats.AlignmentsFound++
							}
						}
					}
				}
				d = !d
			}
		}
		if err = scanner.Err(); err != nil {
			return
		}
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es
	stats.AlignmentTestsPerSec = uint64(float64(stats.AlignmentsTested) / es)
	return
}
