// Select fasta entries from a directory or file containing sequential fastas.
// Pre-process interleaved files with "split" first, and optionally run split
// on the output to split the selected file into smaller parts.

// The search query can be a partial or exact match of any part of the fasta
// description. eg, to select by the UniProtKB entry name suffix:
// -search="_HUMAN" or by the description, -search="Chimpanzee adenovirus Y25"

package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jbrough/blastr"
	"github.com/rs/zerolog/log"
)

func main() {
	search := flag.String("search", "", "search id, name or descriptive text")
	in := flag.String("in", "", "candidate fasta file or directory")
	out := flag.String("out", "", "fasta out file")
	flag.Parse()

	t := time.Now()

	outCh := make(chan [2]string)

	var fastas []string

	if strings.HasSuffix(*in, ".fa") {
		fastas = append(fastas, *in)
	} else {
		if err := filepath.Walk(*in, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".fa") {
				fastas = append(fastas, path)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}

	file, err := os.Create(*out)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	go func() {
		for f := range outCh {
			log.Print("found " + f[0])
			line := f[0] + "\n" + f[1] + "\n"
			if _, err = file.WriteString(line); err != nil {
				return
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(len(fastas))

	for _, path := range fastas {
		go func(path string) {
			defer wg.Done()
			if err := blastr.Select(path, *search, outCh); err != nil {
				panic(err)
			}
			log.Print("finished " + path)
		}(path)
	}

	wg.Wait()

	e := time.Now().Sub(t).Seconds()
	log.Debug().Float64("elapased_secs", e).Send()
}
