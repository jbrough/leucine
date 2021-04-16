package main

import (
	"flag"

	"github.com/jbrough/leucine/runner"
)

// TODO refactor with a cli lib
func main() {

	// import
	src := flag.String("src", "", "source file or directory")
	dst := flag.String("dst", "", "out folder")
	split_after := flag.Int("split", 1e7, "number of entries per fasta file eg 1e7")

	// select
	query := flag.String("search", "", "search id, name or descriptive text")
	jv := flag.Bool("j", false, "output CRLF-delimited JSON objects for streaming results")

	// search
	query_src := flag.String("query", "", "query file")
	candidates_src := flag.String("candidates", "", "candidates file or dir")
	ngram_len := flag.Int("n", 5, "min match length")

	// score
	min_score := flag.Int("min", 0, "min score")
	max_score := flag.Int("max", 9999, "max score")
	filter := flag.String("filter", "", "filter by AA subsequence")

	flag.Parse()

	vals := flag.Args()
	cmd := vals[0]

	var err error
	switch cmd {

	case "import":
		err = runner.Import(*src, *dst, *split_after)

	case "select":
		err = runner.Select(*src, *dst, *query)

	case "search":
		err = runner.Search(*query_src, *candidates_src, *ngram_len, *jv)

	case "score":
		err = runner.Score(*min_score, *max_score, *filter, *jv)

	case "pretty":
		err = runner.Pretty()
	}

	if err != nil {
		panic(err)
	}
}
