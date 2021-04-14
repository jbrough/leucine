package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
	"github.com/jbrough/leucine/search"
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

	outCh := make(chan search.Alignment)

	info := metrics.AlignInfo{*query, *candidates, metrics.AlignStats{}}

	paths, err := io.PathsFromOpt(*candidates)
	if err != nil {
		panic(err)
	}

	go func() {
		for a := range outCh {
			j, err := json.Marshal(a)
			if err != nil {
				panic(err)
			}
			if *jv {
				fmt.Println(string(j))
			}
		}
	}()

	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			stats, err := search.Align(*query, path, *n, outCh)
			if err != nil {
				panic(err)
			}
			stats.FastaFile = path
			fmt.Println(stats.AsJSON())

			info.Stats.Add(stats)
			info.Stats.Stats = append(info.Stats.Stats, stats)
		}(path, &wg)
	}

	wg.Wait()

	j, err := json.Marshal(&info)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
