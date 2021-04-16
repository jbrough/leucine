package fasta

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var faDefResults = [2]string{
	">sp|Q6GZX4|001R_FRG3G Putative transcription factor 001R OS=Frog virus 3 (isolate Goorha) OX=654924 GN=FV3-001R PE=4 SV=1",
	">sp|Q6GZX1|004R_FRG3G Uncharacterized protein 004R OS=Frog virus 3 (isolate Goorha) OX=654924 GN=FV3-004R PE=4 SV=1",
}

var faSeqResults = [2]string{
	"MAFSAEDVLKEYDRRRRMEALLLSLYYPNDRKLLDYKEWSPPRVQVECPKAPVEWNNPPSEKGLIVGHFSGIKYKGEKAQASEVDVNKMCCWVSKFKDAMRRYQGIQTCKIPGKVLSDLDAKIKAYNLTVEGVEGFVRYSRVTKQHVAAFLKELRHSKQYENVNLIHYILTDKRVDIQHLEKDLVKDFKALVESAHRMRQGHMINVKYILYQLLKKHGHGPDGPDILTVKTGSKGVLYDDSFRKIYTDLGWKFTPL",
	"MNAKYDTDQGVGRMLFLGTIGLAVVVGGLMAYGYYYDGKTPSSGTSFHTASPSFSSRYRY",
}

func TestFromFasta(t *testing.T) {
	f, err := os.Open("testdata/interleaved.fa")
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		if err := ParseFasta(scanner, ch); err != nil {
			t.Fatal(err)
		}
	}()

	var tests int
	var i int
	for entry := range ch {
		s := string(entry)
		l := strings.Split(s, "\n")

		got := l[0]
		want := faDefResults[i]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("def mismatch on seq %d (-want +got):\n%s", i, diff)
		}
		tests++

		got = l[1]
		want = faSeqResults[i]
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("/translation mismatch on seq %d (-want +got):\n%s", i, diff)
		}
		tests++
		i++
	}

	if tests != 2 {
		t.Errorf("Expected 4 results, got %d", tests)
	}

}
