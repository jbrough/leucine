package metrics

import "encoding/json"

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
