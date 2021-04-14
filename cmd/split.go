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

	leucine "github.com/jbrough/leucine/pkg"
)

func main() {
	n := flag.Float64("n", 1e7, "number of entries per fasta file eg 1e7")
	in := flag.String("in", "", "fasta source file or directory")
	out := flag.String("out", "", "out folder")
	flag.Parse()

	ts := time.Now()

	info := leucine.SplitInfo{*in, *out, []leucine.SplitStats{}, 0}

	paths, err := leucine.FastaPathsFromOpt(*in)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			var stats leucine.SplitStats
			var err error
			switch ft := filepath.Ext(path); ft {
			case ".seq":
				stats, err = leucine.SplitSequence(path, *out, int(*n))
			case ".fasta", ".fa", ".faa":
				stats, err = leucine.SplitFasta(path, *out, int(*n))
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
