package blastr

import (
	"bufio"
	"bytes"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/twmb/murmur3"
)

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

type Stats struct {
	SequencesSearched    int     `json:"total_sequeneces_searched"`
	AlignmentsFound      int     `json:"alignments_found`
	AlignmentsTested     uint64  `json:"alignments_tested"`
	RuntimeSecs          float64 `json:"runtime_secs"`
	AlignmentTestsPerSec uint64  `json:"alginment_tests_per_sec"`
}

func (s Stats) Add(sa Stats) Stats {
	return Stats{
		s.SequencesSearched + sa.SequencesSearched,
		s.AlignmentsFound + sa.AlignmentsFound,
		s.AlignmentsTested + sa.AlignmentsTested,
		s.RuntimeSecs + sa.RuntimeSecs,
		s.AlignmentTestsPerSec + sa.AlignmentTestsPerSec,
	}
}

type Hit struct {
	QueryId  string `json:"id"`
	QueryPos int    `json:"i"`
	HitId    string `json:"bid"`
	HitPos   int    `json:"bi"`
	HitWord  string `json:"m"`
}

func Align(query_path, test_path string, ngram_n int, out chan Hit) (stats Stats, err error) {
	query_file, err := os.Open(query_path)
	if err != nil {
		return
	}

	query_test := make(map[uint64]bool)
	query_table := make(map[uint64]map[uint64]int)
	query_ids := make(map[uint64]string)

	stats = Stats{}

	d := true
	var id []byte

	scanner := bufio.NewScanner(query_file)
	for scanner.Scan() {
		l := scanner.Bytes()
		if d {
			id = btoid(l)
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
			d = !d
		} else {
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
							out <- Hit{
								QueryId:  query_ids[qid],
								QueryPos: idx,
								HitId:    string(id),
								HitPos:   i,
								HitWord:  string(word),
							}
							skip = i + ngram_n
							stats.AlignmentsFound++
						}
					}
				}
			}
			d = !d
		}
		if err = scanner.Err(); err != nil {
			return
		}
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es
	stats.AlignmentTestsPerSec = stats.AlignmentsTested / uint64(es)
	return
}
