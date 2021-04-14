// Convert an interleaved fasta file into sequential format and optionally split into
// part files containing a max number of entries each.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/jbrough/leucine/fasta"
	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func main() {
	n := flag.Float64("n", 1e7, "number of entries per fasta file eg 1e7")
	in := flag.String("in", "", "source file or directory")
	out := flag.String("out", "", "out folder")
	flag.Parse()

	ts := time.Now()

	info := metrics.SplitInfo{*in, *out, []metrics.SplitStats{}, 0}

	paths, err := io.PathsFromOpt(*in)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			var stats metrics.SplitStats
			var err error
			switch ft := filepath.Ext(path); ft {
			case ".seq":
				stats, err = fasta.FromGenBankSeq(path, *out, int(*n))
			case ".fasta", ".fa", ".faa":
				stats, err = fasta.FromFasta(path, *out, int(*n))
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(stats.AsJSON())

			info.Stats = append(info.Stats, stats)
		}(path, &wg)
	}

	wg.Wait()

	info.RuntimeSecs = time.Now().Sub(ts).Seconds()
	j, err := json.Marshal(&info)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
