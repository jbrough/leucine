// Select fasta entries from a directory or file containing sequential fastas.
// Pre-process interleaved files with "split" first, and optionally run split
// on the output to split the selected file into smaller parts.

// The query query can be a partial or exact match of any part of the fasta
// description. eg, to select by the UniProtKB entry name suffix:
// -query="_HUMAN" or by the description, -search="Chimpanzee adenovirus Y25"

package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jbrough/leucine/fasta"
	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func Select(src, dst, query string) (err error) {
	ts := time.Now()

	info := metrics.SelectInfo{src, dst, query, []metrics.SelectStats{}, 0}

	paths, err := io.PathsFromOpt(src)
	if err != nil {
		panic(err)
	}

	outCh := make(chan [2]string)

	file, err := os.Create(dst)
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
		if strings.Contains(path, dst) {
			continue
		}
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			stats, err := fasta.Select(path, query, outCh)
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
		return
	}
	fmt.Println(string(j))

	return
}
