package runner

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/jbrough/leucine/fasta"
	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func Import(src, dst string, split_after int) (err error) {
	ts := time.Now()
	info := metrics.SplitInfo{src, dst, []metrics.SplitStats{}, 0}

	paths, err := io.PathsFromOpt(src)
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
				stats, err = fasta.FromGenBankSeq(path, dst, int(split_after))
			case ".fasta", ".fa", ".faa":
				stats, err = fasta.FromFasta(path, dst, int(split_after))
			}
			if err != nil {
				return
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

	return
}
