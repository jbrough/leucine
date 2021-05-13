package fasta

import (
	"fmt"
	"strings"
)

type Alignment struct {
	Score   int
	Match   string
	Query   LocalSeq
	Subject LocalSeq
}

type LocalSeq struct {
	Defs []string
	Seq  string
	Idx  []int
}

func (a Alignment) Format() string {
	return fmt.Sprintf(`
    Score: %d
		%s
		-----
		Query  %4d  %s  %4d
		             %s    
		Sbjct  %4d  %s  %4d
		-----
		%s`,
		a.Score, a.Query.Defs,

		a.Query.Idx[0], a.Query.Seq, a.Query.Idx[1],
		a.Match,

		a.Subject.Idx[0], a.Subject.Seq, a.Subject.Idx[1], a.Subject.Defs)
}

func Align(a, b string, word string) Alignment {
	s := 40
	qidx := strings.Index(a, word)
	sidx := strings.Index(b, word)

	if s%2 != 0 {
		s++
	}

	max := (60 - s) / 2

	len_q := len(a)
	len_s := len(b)

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

	qs := a[qx:qy]
	ss := b[sx:sy]

	matches := make(map[int]int)
	var score int
	var prevmatch bool
	var boost int
	var sm string

	for i, r := range qs {
		if qs[i] != ss[i] {
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

	return Alignment{
		Score: score,
		Match: sm,
		Query: LocalSeq{
			Idx: []int{qx, qy},
			Seq: qs,
		},
		Subject: LocalSeq{
			Idx: []int{sx, sy},
			Seq: ss,
		},
	}
}
