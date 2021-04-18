// Format the output from align into a familar BLAST alignment
// This also outputs JSON with the formatted lines of the
// alignment, but it can be pretty-formatted in the console
// by piping to `cmd/pretty.go` or `jq`.

package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jbrough/leucine/search"
)

func Score(min_score, max_score int, filter string, jv bool) (err error) {
	upfilter := strings.ToUpper(filter)

	var test_filter bool
	if filter != "" {
		test_filter = true
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		return
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("cat file.json | prettyjson")
		return
	}

	dec := json.NewDecoder(os.Stdin)

	sidx, err := search.LoadSasas("./sasa/")
	if err != nil {
		return err
	}

	var a search.Alignment
	for {
		if err := dec.Decode(&a); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("decode error %v", err)
		}
		if a.QueryId != "" {
			if test_filter {
				if !strings.Contains(a.Word, upfilter) {
					continue
				}
			}

			score := search.BasicScore(&a, sidx)

			var j []byte
			j, err = json.Marshal(score)
			if err != nil {
				return
			}

			if score.Score >= min_score && score.Score <= max_score {
				if jv {
					fmt.Println(string(j))
				} else {
					fmt.Printf(
						"\n\nScore: %d\n\n\033[34m%s\n\033[37m%s\n%s\n%s\n\nQuery: %s\nSbjct: %s\n",
						score.Score,
						score.Alignment[0],
						//score.Alignment[1],
						score.Alignment[2],
						score.Alignment[3],
						score.Alignment[4],
						score.Query,
						score.Subject,
					)
				}
			}
		}
	}

	return
}
