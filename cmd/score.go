// Format the output from align into a familar BLAST alignment
// This also outputs JSON with the formatted lines of the
// alignment, but it can be pretty-formatted in the console
// by piping to `cmd/pretty.go` or `jq`.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jbrough/leucine"
)

func main() {
	min_score := flag.Int("score", 0, "min score")
	filter := flag.String("filter", "", "filter by AA subsequence")

	flag.Parse()

	var test_filter bool
	if *filter != "" {
		test_filter = true
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("cat file.json | prettyjson")
		return
	}

	dec := json.NewDecoder(os.Stdin)

	var a leucine.Alignment
	for {
		if err := dec.Decode(&a); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("decode error %v", err)
		}
		if a.QueryId != "" {
			if test_filter {
				if !strings.Contains(a.Word, *filter) {
					continue
				}
			}
			score := leucine.BasicScore(&a)
			j, err := json.Marshal(score)
			if err != nil {
				panic(err)
			}

			if score.Score >= *min_score {
				fmt.Println(string(j))
			}
		}

	}
}
