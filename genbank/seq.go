package genbank

import (
	"fmt"
	"strings"
)

type Cds struct {
	CodonStart  string
	Gene        string
	Product     string
	ProteinId   string
	Region      string
	Translation string
}

type Locus struct {
	Accession string
	Version   string
	Organism  string
	Organelle string
	MolType   string
	DbXRef    string
	cds       []Cds
}

func (l *Locus) Fasta() []byte {
	c := l.Cds()
	def := fmt.Sprintf(
		"gb|%s|%s organisim=%q organelle=%q mol_type=%q db_xref=%q gene=%q cds=%q codon_start=%q version=%q dbsource=%q\n",
		c.ProteinId,
		c.Product,
		l.Organism,
		l.Organelle,
		l.MolType,
		l.DbXRef,
		c.Gene,
		c.Region,
		c.CodonStart,
		l.Version,
		l.Accession,
	)

	lines := ">" + strings.Replace(def, ">", "gt", -1) + c.Translation + "\n"
	return []byte(lines)
}

func (l *Locus) NewCds() *Cds {
	l.cds = append(l.cds, Cds{})
	return l.Cds()
}

func (l *Locus) Cds() *Cds {
	return &l.cds[len(l.cds)-1]
}
