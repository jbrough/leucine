package search

import (
	"bufio"
	"os"
	"time"

	"github.com/jbrough/leucine/metrics"
)

type Alignment struct {
	QueryId    string        `json:"qid"`
	QueryIdx   int           `json:"qi"`
	SubjectId  string        `json:"sid"`
	SubjectIdx int           `json:"si"`
	Word       string        `json:"w"`
	QuerySeq   LocalSequence `json:"qs,omitempty"`
	SubjectSeq LocalSequence `json:"ss,omitempty"`
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

func Words(seq []byte, n int) (r [][]byte) {
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

func Align(query_path, test_path string, ngram_n int, out chan<- Alignment) (stats metrics.AlignStats, err error) {
	query_file, err := os.Open(query_path)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(query_file)
	index, err := IndexStream(scanner, ngram_n)
	if err != nil {
		return
	}

	query_file.Close()

	test_file, err := os.Open(test_path)
	if err != nil {
		return
	}
	defer test_file.Close()

	scanner = bufio.NewScanner(test_file)
	return SearchStream(scanner, ngram_n, index, out)
}

func SearchStream(scanner *bufio.Scanner, ngram_n int, index *Index, out chan<- Alignment) (stats metrics.AlignStats, err error) {
	ts := time.Now()

	var def []byte
	d := true

	for scanner.Scan() {
		l := scanner.Bytes()
		var skip int
		if d {
			stats.SequencesSearched++

			def = nil
			def = make([]byte, len(l)-1) // remove leading '>'
			copy(def, l[1:])

			d = !d
		} else {
			// dont check self
			if _, ok := index.Match[index.Hash(def)]; !ok {
				for i, word := range Words(l, ngram_n) {
					if skip > 0 && i < skip {
						continue
					} else {
						skip = 0
					}
					stats.AlignmentsTested++
					if _, ok := index.Test(word); ok {
						for qid, tbl := range index.Match {
							if idxs, ok := tbl[index.Hash(word)]; ok {
								skip = i + ngram_n

								id := string(def)
								tmp := make([]byte, len(l))
								copy(tmp, l)

								// TODO: woops
								idx := idxs[0]

								qseq, sseq := localSequences(index.GetRef[qid][1], tmp, idx, i, ngram_n)

								out <- Alignment{
									QueryId:    index.GetKey[qid],
									QueryIdx:   idx,
									SubjectId:  id,
									SubjectIdx: i,
									Word:       string(word),
									QuerySeq:   qseq,
									SubjectSeq: sseq,
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
