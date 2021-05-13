package fasta_test

import (
	"fmt"
	"strings"
	"testing"
)

var def = `>pdb|2HM3|2HM3_1|NEMATOCYST OUTER WALL ANTIGEN, CYSTEINE RICH DOMAIN NW1 mol="NEMATOCYST OUTER WALL ANTIGEN" frag="FIRST CYSTEINE RICH DOMAIN" org="HYDRA VULGARIS" taxid=6087 dbref=Q8IT70:8IT70_HYDAT chain=A start=7 end=31`

func parse(s string) {
	m := strings.Split(s, "|")
	//db := m[0][1:]
	ent := m[1]
	mol := strings.Split(m[2], "_")[1]
	i := strings.Index(s, "chain=")
	c := s[i+6 : i+7]
	fmt.Printf("%s,%s,%s,%s\n", ent, mol, c, "1.5")
}

func TestPdbEndFromFasta(t *testing.T) {

	parse(def)

}
