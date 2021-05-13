package fasta

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/jbrough/leucine/pdb"
)

func ParsePdbEnt(scanner *bufio.Scanner) (ent *pdb.Ent, err error) {
	val := func(l string, i int) string {
		return strings.TrimSuffix(strings.TrimSpace(l[i:]), ";")
	}

	mols := make(map[string]*pdb.Mol)
	var molid string

	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, "HEADER") {
			ent = &pdb.Ent{}
			ent.Date = val(l, 50)[:10]
			ent.Id = val(l, 62)
			continue
		}

		if strings.HasPrefix(l, "TITLE") {
			if ent.Title != "" {
				ent.Title += " "
			}
			ent.Title += val(l, 8)
			continue
		}

		if strings.HasPrefix(l, "COMPND") {
			if strings.HasPrefix(l[10:], "MOL_ID") || strings.HasPrefix(l[11:], "MOL_ID") {
				molid = val(l, 18)
				mols[molid] = &pdb.Mol{Id: molid}

				continue
			}

			if strings.HasPrefix(l[11:], "MOLECULE") {
				mols[molid].Name = val(l, 21)
				continue
			}

			if strings.HasPrefix(l[11:], "FRAGMENT") {
				mols[molid].Fragment = val(l, 21)
				continue
			}

			if strings.HasPrefix(l[11:], "SYNONYM") {
				mols[molid].Synonym = val(l, 20)
				continue
			}

			if strings.HasPrefix(l[11:], "CHAIN") {
				ids := strings.Split(val(l, 18), ",")
				for _, id := range ids {
					mols[molid].Chains = append(mols[molid].Chains, &pdb.Chain{Id: strings.TrimSpace(id)})
				}
				continue
			}

			continue
		}

		if strings.HasPrefix(l, "SOURCE") {
			if strings.HasPrefix(l[10:], "MOL_ID") || strings.HasPrefix(l[11:], "MOL_ID") {
				molid = val(l, 18)

				continue
			}

			if strings.HasPrefix(l[11:], "ORGANISM_SCIENTIFIC") {
				mols[molid].Organism = val(l, 32)
				continue
			}

			if strings.HasPrefix(l[11:], "ORGANISM_TAXID") {
				mols[molid].TaxId = strings.Split(val(l, 27), ",")[0]
				continue
			}
		}

		if strings.HasPrefix(l, "DBREF") {
			id := l[12:13]
			for im, m := range mols {
				for ic, c := range m.Chains {
					if c.Id == id {
						refa := val(l[33:42], 0)
						refb := val(l[43:57], 0)
						c = mols[im].Chains[ic]
						c.DbRef = fmt.Sprintf("%s:%s", refa, refb)
						//mols[mid].Chains[ic] = c
					}
				}
			}
			continue
		}

		if strings.HasPrefix(l, "ATOM") {
			var ms []*pdb.Mol
			for _, mol := range mols {
				ms = append(ms, mol)
			}

			ent.Mols = ms

			break
		}
	}

	return
}
