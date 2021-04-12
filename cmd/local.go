// Format the output from align into a familar BLAST alignment
// This also outputs JSON with the formatted lines of the
// alignment, but it can be pretty-formatted in the console
// by piping to `cmd/pretty.go` or `jq`.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jbrough/leucine"
)

type PrettyAlignment struct {
	Query     string    `json:"query"`
	Subject   string    `json:"sbjct"`
	Alignment [3]string `json:"align"`
}

func main() {
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
			la := a.QuerySeq
			lb := a.SubjectSeq
			lc := ""

			for i, c := range la.A {
				if byte(c) == lb.A[i] {
					lc += string(c)
				} else {
					lc += " "
				}
			}

			strs := [3]string{
				fmt.Sprintf("Query  %3d  %s  %3d", la.X, la.A, la.Y),
				fmt.Sprintf("            %s     ", lc),
				fmt.Sprintf("Sbjct  %3d  %s  %3d", lb.X, lb.A, lb.Y),
			}
			j, err := json.Marshal(&PrettyAlignment{a.QueryName, a.SubjectName, strs})
			if err != nil {
				panic(err)
			}

			fmt.Println(string(j))
		}

	}
}
