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
	Query       string     `json:"query"`
	Candidates  string     `json:"candidates";`
	Stats       AlignStats `json:"stats"`
	RuntimeSecs float64    `json:"runtime_secs"`
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

func (s AlignStats) Add(sa AlignStats) AlignStats {
	return AlignStats{
		SequencesSearched:    s.SequencesSearched + sa.SequencesSearched,
		AlignmentsFound:      s.AlignmentsFound + sa.AlignmentsFound,
		AlignmentsTested:     s.AlignmentsTested + sa.AlignmentsTested,
		AlignmentTestsPerSec: s.AlignmentTestsPerSec + sa.AlignmentTestsPerSec,
	}
}

func (s AlignStats) AsJSON() string {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(j)
}

type Alignment struct {
	QueryId      string `json:"qid"`
	QueryIdx     int    `json:"qi"`
	CandidateId  string `json:"cid"`
	CandidateIdx int    `json:"ci"`
	Word         string `json:"w"`
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

	stats = AlignStats{}

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
							out <- Alignment{
								QueryId:      query_ids[qid],
								QueryIdx:     idx,
								CandidateId:  string(id),
								CandidateIdx: i,
								Word:         string(word),
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
