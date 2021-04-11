package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jbrough/blastr"
	"github.com/rs/zerolog/log"
)

func main() {
	query := flag.String("query", "", "path to queryfasta")
	candidates := flag.String("candidates", "", "path to candidate fasta or directory")
	out := flag.String("out", "", "out file")
	n := flag.Int("n", 6, "length of match")

	_ = out
	flag.Parse()

	t := time.Now()

	outCh := make(chan blastr.Hit)

	if strings.HasSuffix(*candidates, ".fa") {
		fastas = append(fastas, *candidates)
	} else {
		var fastas []string
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
			fmt.Println(string(j))
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(fastas))

	for _, path := range fastas {
		go func(path string) {
			defer wg.Done()
			if err := blastr.Align(*query, path, *n, outCh); err != nil {
				panic(err)
			}
			log.Print("finished " + path)
		}(path)
	}

	wg.Wait()

	e := time.Now().Sub(t).Seconds()
	log.Debug().Float64("elapased_secs", e).Send()
}
