// Select fasta entries from a directory or file containing sequential fastas.
// Pre-process interleaved files with "split" first, and optionally run split
// on the output to split the selected file into smaller parts.

// The query query can be a partial or exact match of any part of the fasta
// description. eg, to select by the UniProtKB entry name suffix:
// -query="_HUMAN" or by the description, -search="Chimpanzee adenovirus Y25"

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jbrough/leucine/fasta"
	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func main() {
	query := flag.String("search", "", "search id, name or descriptive text")
	in := flag.String("in", "", "candidate fasta file or directory")
	out := flag.String("out", "", "fasta out file")
	flag.Parse()

	ts := time.Now()

	info := metrics.SelectInfo{*in, *out, *query, []metrics.SelectStats{}, 0}

	paths, err := io.PathsFromOpt(*in)
	if err != nil {
		panic(err)
	}

	outCh := make(chan [2]string)

	file, err := os.Create(*out)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	go func() {
		for f := range outCh {
			line := f[0] + "\n" + f[1] + "\n"
			if _, err = file.WriteString(line); err != nil {
				return
			}
		}
	}()

	wg := sync.WaitGroup{}

	for _, path := range paths {
		if strings.Contains(path, *out) {
			continue
		}
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			stats, err := fasta.Select(path, *query, outCh)
			if err != nil {
				panic(err)
			}
			info.Stats = append(info.Stats, stats)
			j, err := json.Marshal(info)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(j))

		}(path)
	}

	wg.Wait()

	info.RuntimeSecs = time.Now().Sub(ts).Seconds()

	j, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))

}
