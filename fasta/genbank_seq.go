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

func parseTranslationLn(b []byte) ([]byte, bool) {
	li := len(b) - 1
	if b[li] == '"' {
		return b[:li], true
	}

	return b, false
}

func FromGenBankSeq(in, out string, limit int) (stats metrics.SplitStats, err error) {
	ts := time.Now()

	stats.Source = in
	name := strings.TrimSuffix(filepath.Base(in), "")

	file, err := os.Open(in)
	if err != nil {
		return
	}
	defer file.Close()

	files, err := io.NewPartFiles(out, name)
	if err != nil {
		return
	}

	stats.Splits = append(stats.Splits, files.Name())

	scanner := bufio.NewScanner(file)

	locus := &genbank.Locus{}
	var cds *genbank.Cds
	var inseq bool
	var incds bool
	var count int
	for scanner.Scan() {
		l := scanner.Bytes()

		if bytes.HasPrefix(l, []byte("LOCUS")) {
			locus = &genbank.Locus{}
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

			if bytes.HasPrefix(l[21:], []byte("/translation=")) {
				inseq = true
				b, last := parseTranslationLn(l[35:])
				cds.Translation += string(b)
				if last {
					inseq = false
				}
			}

			if inseq {
				b, last := parseTranslationLn(l[21:])
				cds.Translation += string(b)
				if last {
					inseq = false

					if count > 0 && count%limit == 0 {
						if err = files.Cycle(); err != nil {
							return
						}
						stats.Splits = append(stats.Splits, files.Name())
					}

					b := locus.CdsBytes()
					if err = files.Write(b); err != nil {
						return
					}
					count++
				}
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	if err = files.Close(); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}
