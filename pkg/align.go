package leucine

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
	QueryId     string        `json:"qid"`
	QueryIdx    int           `json:"qi"`
	SubjectId   string        `json:"sid"`
	SubjectIdx  int           `json:"si"`
	Word        string        `json:"w"`
	QuerySeq    LocalSequence `json:"qs,omitempty"`
	SubjectSeq  LocalSequence `json:"ss,omitempty"`
	QueryName   string        `json:"qn,omitempty"`
	SubjectName string        `json:"sn,omitempty"`
}

type LocalSequence struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	A string `json:"a"`
}

func localSequences(qseq, sseq []byte, qidx, sidx, s int) (LocalSequence, LocalSequence) {
	if s%2 != 0 {
		s++
	}

	max := (60 - s) / 2

	len_q := len(qseq)
	len_s := len(sseq)

	len_qb := len_q - qidx
	len_sb := len_s - sidx

	var x, y int

	if max < qidx && max < sidx {
		x = max
	} else {
		if qidx < sidx {
			x = qidx
		} else {
			x = sidx
		}
	}

	if max < len_qb && max < len_sb {
		y = max
	} else {
		if len_qb < len_sb {
			y = len_qb
		} else {
			y = len_sb
		}
	}

	qx := qidx - x
	qy := qidx + y
	sx := sidx - x
	sy := sidx + y

	var qsuffix, ssuffix string
	if qy == len_q+1 {
		qsuffix = "*"
	}
	if sy == len_s+1 {
		ssuffix = "*"
	}

	qs := string(qseq[qx:qy])
	ss := string(sseq[sx:sy])

	return LocalSequence{
			X: qx,
			Y: qy,
			A: string(qs) + qsuffix,
		}, LocalSequence{
			X: sx,
			Y: sy,
			A: string(ss) + ssuffix,
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
			if bytes.IndexByte(l, '>') == -1 {
				panic("TODO: Interleaved fastas are not currently supported. Please pre-process with `split`")
			}
			d = !d
			id = btoid(l)

			// TODO: go copies slices by value but the value is a reference to a
			// backing array. There was a race here that meant b sometimes referenced
			// the next line. This fixes it but I need to revist this implementation.
			tmp := make([]byte, len(l))
			copy(tmp, l)
			b = tmp

		} else {
			query_index := make(map[uint64]int)
			for i, word := range words(l, ngram_n) {
				h := hash(word)
				query_test[h] = true
				query_index[h] = i
			}
			query_table[hash(id)] = query_index
			query_ids[hash(id)] = string(id)

			tmp := make([]byte, len(l))
			copy(tmp, l)

			query_detail[hash(id)] = [2][]byte{b, tmp}

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

			tmp := make([]byte, len(l))
			copy(tmp, l)
			b = tmp

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

								tmp := make([]byte, len(l))
								copy(tmp, l)

								qseq, sseq := localSequences(query_detail[qid][1], tmp, idx, i, ngram_n)
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
