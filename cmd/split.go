// Convert an interleaved fasta file into sequential format and optionally split into
// part files containing a max number of entries each.

package main

import (
	"flag"

	"github.com/jbrough/blastr"
)

func main() {
	entries := flag.Float64("entries", 1e7, "number of entries per fasta file eg 1e7")
	in := flag.String("in", "", "fasta source file")
	out := flag.String("out", "", "out folder")
	flag.Parse()

	limit := int(*entries)

	if err := blastr.SplitFasta(*in, *out, limit); err != nil {
		panic(err)
	}
}
