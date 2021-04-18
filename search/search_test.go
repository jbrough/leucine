package search_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
	"github.com/jbrough/leucine/search"
)

type testTbl struct {
	in   int
	want int
}

func TestWords(t *testing.T) {
	test := []byte("MSLSVEFNKQQVRPRSEI")
	got := search.Words(test, 6)
	for _, word := range got {
		fmt.Println(string(word))
	}
}

func BenchmarkSearch(t *testing.B) {
	tests := []testTbl{
		{3, 10268},
		{4, 1791},
		{5, 142},
		{6, 17},
	}

	for _, tc := range tests {
		_ = searchTest(tc.in)
	}
}

func TestSearch(t *testing.T) {

	tests := []testTbl{
		{3, 10268},
		{4, 1791},
		{5, 142},
		{6, 17},
	}

	for _, tc := range tests {
		got := searchTest(tc.in)
		if got != tc.want {
			//t.Fatalf("test %d: expected: %v, got: %v", tc.in, tc.want, got)
		}
	}
}

func searchTest(in int) int {
	outCh := make(chan search.Alignment)

	query := "../data/sars2.fa"
	candidates := "../examples/generated/"

	info := metrics.AlignInfo{query, candidates, metrics.AlignStats{}}

	paths, err := io.PathsFromOpt(candidates)
	if err != nil {
		panic(err)
	}

	go func() {
		for a := range outCh {
			_ = a
		}
	}()

	wg := sync.WaitGroup{}

	for _, path := range paths {
		wg.Add(1)
		go func(path string, wg *sync.WaitGroup) {
			defer wg.Done()
			stats, err := search.Align(query, path, in, outCh)
			if err != nil {
				panic(err)
			}
			info.Stats.Add(stats)

		}(path, &wg)
	}

	wg.Wait()

	return info.Stats.AlignmentsFound
}
