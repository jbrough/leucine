package fasta

import (
	"bufio"
	"bytes"
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

	parseTranslationLn := func(b []byte) ([]byte, bool) {
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
		l := scanner.Bytes()

		if bytes.HasPrefix(l, []byte("LOCUS")) {
			locus = genbank.Locus{}
			continue
		}

		if bytes.HasPrefix(l, []byte("ACCESSION")) {
			locus.Accession = string(l[12:])
			continue
		}

		if bytes.HasPrefix(l, []byte("VERSION")) {
			locus.Version = string(l[12:])
			continue
		}

		s := len(l)

		if incds && s > 5 {
			if l[5] != ' ' {
				incds = false
			}
		}
		if s > 21 {

			if bytes.HasPrefix(l[21:], []byte("/organism")) {
				locus.Organism = string(l[32 : s-1])
				continue
			}

			if bytes.HasPrefix(l[21:], []byte("/organelle")) {
				locus.Organelle = string(l[33 : s-1])
				continue
			}

			if bytes.HasPrefix(l[21:], []byte("/mol_type")) {
				locus.MolType = string(l[32 : s-1])
				continue
			}

			if bytes.HasPrefix(l[21:], []byte("/db_xref")) {
				locus.DbXRef = string(l[31 : s-1])
				continue
			}

			if bytes.HasPrefix(l[5:], []byte("CDS")) {
				cds = locus.NewCds()
				cds.Region = string(l[21:])
				incds = true
				continue
			}

			if incds && bytes.HasPrefix(l[21:], []byte("/gene=")) {
				cds.Gene = string(l[28 : s-1])
				continue
			}

			if incds && bytes.HasPrefix(l[21:], []byte("/codon_start=")) {
				cds.CodonStart = string(l[34:])
				continue
			}

			if incds && bytes.HasPrefix(l[21:], []byte("/product=")) {
				cds.Product = string(l[31 : s-1])
				continue
			}

			if incds && bytes.HasPrefix(l[21:], []byte("/protein_id=")) {
				cds.ProteinId = string(l[34 : s-1])
				continue
			}

			if inseq {
				b, last := parseTranslationLn(l[21:])
				seq += string(b)
				if last {
					cds.Translation = seq
					inseq = false
					seq = ""

					entries <- locus.CdsBytes()
				}
			} else if bytes.HasPrefix(l[21:], []byte("/translation=")) {
				inseq = true
				b, last := parseTranslationLn(l[35:])
				seq += string(b)
				if last {
					cds.Translation = seq
					inseq = false
					seq = ""

					entries <- locus.CdsBytes()
				}
			}
		}

	}

	return scanner.Err()
}
