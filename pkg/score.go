package leucine

import (
	"fmt"
)

type Score struct {
	Query     string    `json:"query"`
	Subject   string    `json:"sbjct"`
	Alignment [3]string `json:"align"`
	Score     int       `json:"score"`
}

func BasicScore(a *Alignment) *Score {
	la := a.QuerySeq
	lb := a.SubjectSeq
	var sm string

	n := 0
	for i, r := range la.A {
		if la.A[i] != lb.A[i] {
			sm += " "
		} else {
			sm += string(r)
			n++
		}
	}

	strs := [3]string{
		fmt.Sprintf("Query  %4d  %s  %4d", la.X, la.A, la.Y),
		fmt.Sprintf("             %s      ", sm),
		fmt.Sprintf("Sbjct  %4d  %s  %4d", lb.X, lb.A, lb.Y),
	}

	return &Score{a.QueryName, a.SubjectName, strs, n}
}
