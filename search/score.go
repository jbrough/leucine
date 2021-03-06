package search

import (
	"fmt"
)

type Score struct {
	Query     string    `json:"query"`
	Subject   string    `json:"sbjct"`
	Alignment [5]string `json:"align"`
	Score     int       `json:"score"`
}

func BasicScore(a *Alignment) *Score {
	la := a.QuerySeq
	lb := a.SubjectSeq
	var sm string

	// These subsequences are centered on a word of n-length as specified in the
	// alignment search.
	// Individual matches between the sequences are to be expected but should be
	// considered arbitrary.
	// Individual matches therefore score 0 unless there is a subsquent match
	// with a gap of not more than 1

	matches := make(map[int]int)
	var score int
	var prevmatch bool
	var boost int

	for i, r := range la.A {
		if la.A[i] != lb.A[i] {
			if prevmatch {
				matches[i-1] = 1
			}
			prevmatch = false
			sm += " "
		} else {
			if _, ok := matches[i-2]; ok {
				score++
			}

			if prevmatch {
				boost += 1
				score += boost
			} else {
				prevmatch = true
				boost = 5
			}
			sm += string(r)
		}
	}

	strs := [5]string{
		fmt.Sprintf(" SASA        %s     ", ""),
		fmt.Sprintf("             %s     ", ""),
		fmt.Sprintf("Query  %4d  %s  %4d", la.X, la.A, la.Y),
		fmt.Sprintf("             %s      ", sm),
		fmt.Sprintf("Sbjct  %4d  %s  %4d", lb.X, lb.A, lb.Y),
	}

	return &Score{a.QueryId, a.SubjectId, strs, score}
}
