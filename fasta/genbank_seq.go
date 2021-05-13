package fasta

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jbrough/leucine/genbank"
	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func FromGenBankSeq(in, out string, limit int) (stats metrics.SplitStats, err error) {
	ts := time.Now()

	stats.Source = in
	name := strings.TrimSuffix(filepath.Base(in), "")

	file, err := os.Open(in)
	if err != nil {
		return
	}
	defer file.Close()

	files, err := io.NewPartFiles(out, name, limit)
	if err != nil {
		return
	}
	defer files.Close()

	scanner := bufio.NewScanner(file)
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		for entry := range ch {
			part, newpart, err := files.Write(entry)
			if err != nil {
				return
			}

			if newpart {
				stats.Splits = append(stats.Splits, part)
			}
		}
	}()

	if err = ParseGenBankSeq(scanner, ch); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}

func ParseGenBankSeq(scanner *bufio.Scanner, entries chan<- []byte) (err error) {
	parseTranslationLn := func(b string) (string, bool) {
		li := len(b) - 1
		if b[li] == '"' {
			return b[:li], true
		}

		return b, false
	}

	var locus genbank.Locus
	var cds *genbank.Cds

	var inseq bool
	var incds bool
	var seq string

	for scanner.Scan() {
		l := scanner.Text()

		if strings.HasPrefix(l, "LOCUS") {
			locus = genbank.Locus{}
			continue
		}

		if strings.HasPrefix(l, "ACCESSION") {
			locus.Accession = l[12:]
			continue
		}

		if strings.HasPrefix(l, "VERSION") {
			locus.Version = l[12:]
			continue
		}

		s := len(l)

		if incds && s > 5 {
			if l[5] != ' ' {
				incds = false
			}
		}
		if s > 21 {

			if strings.HasPrefix(l[21:], "/organism") {
				locus.Organism = l[32 : s-1]
				continue
			}

			if strings.HasPrefix(l[21:], "/organelle") {
				locus.Organelle = l[33 : s-1]
				continue
			}

			if strings.HasPrefix(l[21:], "/mol_type") {
				locus.MolType = l[32 : s-1]
				continue
			}

			if strings.HasPrefix(l[21:], "/db_xref") {
				locus.DbXRef = l[31 : s-1]
				continue
			}

			if strings.HasPrefix(l[5:], "CDS") {
				cds = locus.NewCds()
				cds.Region = l[21:]
				incds = true
				continue
			}

			if incds && strings.HasPrefix(l[21:], "/gene=") {
				cds.Gene = l[28 : s-1]
				continue
			}

			if incds && strings.HasPrefix(l[21:], "/codon_start=") {
				cds.CodonStart = l[34:]
				continue
			}

			if incds && strings.HasPrefix(l[21:], "/product=") {
				cds.Product = l[31 : s-1]
				continue
			}

			if incds && strings.HasPrefix(l[21:], "/protein_id=") {
				cds.ProteinId = l[34 : s-1]
				continue
			}

			if inseq {
				b, last := parseTranslationLn(l[21:])
				seq += b
				if last {
					cds.Translation = seq
					inseq = false
					seq = ""

					entries <- locus.Fasta()
				}
			} else if strings.HasPrefix(l[21:], "/translation=") {
				inseq = true
				b, last := parseTranslationLn(l[35:])
				seq += b
				if last {
					cds.Translation = seq
					inseq = false
					seq = ""

					entries <- locus.Fasta()
				}
			}
		}

	}

	return scanner.Err()
}
