package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jbrough/blastr"
)

func main() {
	query := flag.String("query", "", "query fasta path")
	candidates := flag.String("candidates", "", "candidate fasta path or directory path")
	out := flag.String("out", "", "out file")
	n := flag.Int("n", 6, "length of match")
	jv := flag.Bool("j", false, "output CRLF-delimited JSON objects for streaming results")
	v := flag.Bool("v", false, "enable JSON logging")

	_ = v
	_ = out
	flag.Parse()

	outCh := make(chan blastr.Hit)

	info := blastr.AlignInfo{*query, *candidates, blastr.Stats{}, 0}

	var fastas []string

	if strings.HasSuffix(*candidates, ".fa") {
		fastas = append(fastas, *candidates)
	} else {
		if err := filepath.Walk(*candidates, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".fa") {
				fastas = append(fastas, path)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}

	go func() {
		for hit := range outCh {
			j, err := json.Marshal(hit)
			if err != nil {
				panic(err)
			}
			if *jv {
				fmt.Println(string(j))
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(fastas))

	for _, path := range fastas {
		go func(path string) {
			defer wg.Done()
			stats, err := blastr.Align(*query, path, *n, outCh)
			if err != nil {
				panic(err)
			}
			stats.FastaFile = path
			fmt.Println(stats.asJSON())

			info.stats.Add(stats)
		}(path)
	}

	wg.Wait()

	j, err := json.Marshal(&info)
	if err != nil {
		panic(err)
	}

	fmt.Prinln(string(j))
}
