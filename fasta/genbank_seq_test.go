package fasta_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jbrough/leucine/fasta"
)

var gbDefResults = [4]string{
	`>gb|CAH72364.1|family with sequence similarity 72, member A organisim="Homo sapiens" organelle="" mol_type="genomic DNA" db_xref="UniProtKB/TrEMBL:Q5TYM6" gene="FAM72A" cds="join(64628..64779,68739..68816,72747..72871,81250..81344)" codon_start="1" version="CR407567.2" dbsource="CR407567"`,
	`>gb|CAH72365.1|family with sequence similarity 72, member A organisim="Homo sapiens" organelle="" mol_type="genomic DNA" db_xref="UniProtKB/Swiss-Prot:Q5TYM5" gene="FAM72A" cds="join(66598..66749,68739..68816,72747..72871,81250..81344)" codon_start="1" version="CR407567.2" dbsource="CR407567"`,
	`>gb|CAH72366.1|family with sequence similarity 72, member A organisim="Homo sapiens" organelle="" mol_type="genomic DNA" db_xref="UniProtKB/Swiss-Prot:Q5TYM5" gene="FAM72A" cds="join(66598..66629,68739..68816,72747..72871,81250..81344)" codon_start="1" version="CR407567.2" dbsource="CR407567"`,
	`>gb|CAH72364.1|family with sequence similarity 72, member A organisim="Homo sapiens" organelle="" mol_type="genomic DNA" db_xref="UniProtKB/TrEMBL:Q5TYM6" gene="FAM72A" cds="join(64628..64779,68739..68816,72747..72871,81250..81344)" codon_start="1" version="" dbsource=""`,
}

var gbSeqResults = [4]string{
	"MPTTTALRWTAAARVRRKGRGGWVPAALRSVSQDGVPGCTVMGGETRSPENAVDFTGRCYFTKICKCKLKDIACLKCGNIVGYHVIVPCSSCLLSCNNGHFWMFHSQAVYDINRLDSTGVNVLLWGNLPEIEESTDEDVLNISAEECIR",
	"MSTNICSFKDRCVSILCCKFCKQVLSSRGMKAVLLADTEIDLFSTDIPPTNAVDFTGRCYFTKICKCKLKDIACLKCGNIVGYHVIVPCSSCLLSCNNGHFWMFHSQAVYDINRLDSTGVNVLLWGNLPEIEESTDEDVLNISAEECIR",
	"MSTNICSFKDSAVDFTGRCYFTKICKCKLKDIACLKCGNIVGYHVIVPCSSCLLSCNNGHFWMFHSQAVYDINRLDSTGVNVLLWGNLPEIEESTDEDVLNISAEECIR",
	"NEXTLOCUS",
}

func TestFromGenBankSeq(t *testing.T) {
	f, err := os.Open("testdata/gb_multi_cds.seq")
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		if err := fasta.ParseGenBankSeq(scanner, ch); err != nil {
			panic(err)
		}
	}()

	var tests int
	var i int
	for entry := range ch {
		s := string(entry)
		l := strings.Split(s, "\n")

		got := l[0]
		want := gbDefResults[i]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("def mismatch on def %d (-want +got):\n%s", i, diff)
		}
		tests++

		got = l[1]
		want = gbSeqResults[i]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("/translation mismatch on CDS %d (-want +got):\n%s", i, diff)
		}
		tests++
		i++
	}

	if tests != 8 {
		t.Errorf("Expected 4 results, got %d", tests)
	}
}
