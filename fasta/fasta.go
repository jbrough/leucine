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

	files, err := io.NewPartFiles(out, name)
	if err != nil {
		return
	}

	stats.Splits = append(stats.Splits, files.Name())

	scanner := bufio.NewScanner(file)
	var count int
	for scanner.Scan() {
		l := scanner.Bytes()

		header := bytes.IndexByte(l, '>') == 0

		if count != 0 && header {
			if err = files.NewLine(); err != nil {
				return
			}
		}

		if header && count > 0 && count%limit == 0 {
			if err = files.Cycle(); err != nil {
				return
			}
			stats.Splits = append(stats.Splits, files.Name())
		}

		if err = files.Write(l); err != nil {
			return
		}

		if header {
			count++
			if err = files.NewLine(); err != nil {
				return
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	if err = files.Close(); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}
