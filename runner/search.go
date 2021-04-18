package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
	"github.com/jbrough/leucine/search"
)

func Search(query_src, candidates_src string, ngram_len int, jv bool) (err error) {
	ts := time.Now()
	outCh := make(chan search.Alignment)

	info := metrics.AlignInfo{query_src, candidates_src, metrics.AlignStats{}}

	paths, err := io.PathsFromOpt(candidates_src)
	if err != nil {
		return
	}

	query_file, err := os.Open(query_src)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(query_file)
	index, err := search.IndexStream(scanner, ngram_len)
	if err != nil {
		return
	}

	query_file.Close()

	go func() {
		for a := range outCh {
			j, err := json.Marshal(a)
			if err != nil {
				return
			}
			if jv {
				fmt.Println(string(j))
			}
		}
	}()

	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			stats, err := search.Align(index, path, ngram_len, outCh)
			if err != nil {
				return
			}
			stats.FastaFile = path
			fmt.Println(stats.AsJSON())

			info.Stats.Add(stats)
			info.Stats.Stats = append(info.Stats.Stats, stats)
		}(path, &wg)
	}

	wg.Wait()

	es := time.Now().Sub(ts).Seconds()
	info.Stats.RuntimeSecs = es

	j, err := json.Marshal(&info)
	if err != nil {
		return
	}

	fmt.Println(string(j))

	return
}
