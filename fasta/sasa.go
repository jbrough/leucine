package fasta

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/jbrough/leucine/lib"
	"github.com/jbrough/leucine/pdb"
)

func ParseSasaOut(scanner *bufio.Scanner) (chains []*pdb.Chain, err error) {
	aa := lib.AAA2AA()

	var lastchain string
	var seq string
	var score []float64
	var idx []int

	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "#") {
			continue
		}
		if len(s) == 0 { // next is EOF
			continue
		}

		chain := s[4:5]
		if len(idx) > 0 && chain != lastchain {
			sm := make(map[string][]float64)
			sm["default"] = score
			chains = append(chains, &pdb.Chain{
				Id: lastchain, Seq: seq, Asa: sm, Start: idx[0], End: idx[len(idx)-1],
			})

			seq = ""
			score = nil
			idx = nil
		}

		var pos int

		pos, err = strconv.Atoi(strings.TrimSpace(s[6:11]))
		if err != nil {
			continue // TODO: what's meaning of these alternative residues
		}

		idx = append(idx, pos)

		if chain == lastchain && pos != idx[len(idx)-2]+1 {
			for i := idx[len(idx)-2] + 1; i < pos; i++ {
				seq += "-"
				score = append(score, -1)
			}
		}

		r := strings.TrimSpace(s[12:15])
		a, ok := aa[r]
		if !ok {
			//err = errors.New("Residue not found: " + r)
			continue
		}

		var sc float64
		sc, err = strconv.ParseFloat(strings.TrimSpace(string(s[20:25])), 64)
		if err != nil {
			return
		}

		seq += a
		score = append(score, sc)
		lastchain = chain
	}

	if err = scanner.Err(); err != nil {
		return
	}

	if len(idx) == 0 {
		return
	}

	sm := make(map[string][]float64)
	sm["default"] = score
	chains = append(chains, &pdb.Chain{
		Id: lastchain, Seq: seq, Asa: sm, Start: idx[0], End: idx[len(idx)-1],
	})

	return
}
