package runner

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
	"github.com/jbrough/leucine/search"
)

func Search(query, candidates string, ngram_len int, jv bool) (err error) {
	outCh := make(chan search.Alignment)

	info := metrics.AlignInfo{query, candidates, metrics.AlignStats{}}

	paths, err := io.PathsFromOpt(candidates)
	if err != nil {
		return
	}

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
			stats, err := search.Align(query, path, ngram_len, outCh)
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

	j, err := json.Marshal(&info)
	if err != nil {
		return
	}

	fmt.Println(string(j))

	return
}
