package fasta

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
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

	defer f.Close()

	fi, err := os.Open("testdata/interleaved.fa")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()

	scanner := bufio.NewScanner(f)
	ch := make(chan Entry)

	go func() {
		defer close(ch)
		if err := Parse(scanner, ch); err != nil {
			panic(err)
		}
	}()

	for e := range ch {
		fmt.Println(e)
		if _, err = fi.Seek(e.Src.ByteOffset, 0); err != nil {
			panic(err)
		}
		b := make([]byte, e.Src.Bytes)
		if _, err = io.ReadAtLeast(fi, b, int(e.Src.Bytes)); err != nil {
			panic(err)
		}

		fmt.Println(string(b))
		fmt.Println("----")

	}

	t.Fail()
}
