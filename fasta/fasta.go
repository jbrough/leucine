package fasta

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func FromFasta(in, out string, limit int) (stats metrics.SplitStats, err error) {
	ts := time.Now()

	stats.Source = in
	name := strings.TrimSuffix(filepath.Base(in), ".fasta")

	file, err := os.Open(in)
	if err != nil {
		return
	}
	defer file.Close()

	files, err := io.NewPartFiles(out, name, limit)
	if err != nil {
		return
	}
	defer files.Close()

	stats.Splits = append(stats.Splits, files.Name())

	scanner := bufio.NewScanner(file)
	ch := make(chan []byte)

	go func() {
		defer close(ch)
		for entry := range ch {
			part, newpart, err := files.Write(entry)
			if err != nil {
				return
			}

			if newpart {
				stats.Splits = append(stats.Splits, part)
			}
		}
	}()

	if err = ParseFasta(scanner, ch); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}

func ParseFasta(scanner *bufio.Scanner, entries chan<- []byte) (err error) {
	var entry []byte

	for scanner.Scan() {
		l := scanner.Bytes()

		def := bytes.IndexByte(l, '>') == 0

		if def {
			if entry == nil { // first line of file
				entry = append(entry, l...)
				entry = append(entry, '\n')
			} else { // back on a definition line with a def and seq in the tmp var
				entry = append(entry, '\n')
				entries <- entry
				entry = nil

				entry = append(entry, l...)
				entry = append(entry, '\n')
			}
		} else {
			entry = append(entry, l...)
		}
	}

	return scanner.Err()
}
