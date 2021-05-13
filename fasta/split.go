package fasta

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/metrics"
)

func Split(in, out string, limit int) (stats metrics.SplitStats, err error) {
	ts := time.Now()

	stats.Source = in
	name := strings.TrimSuffix(filepath.Base(in), ".fasta")

	files, err := io.NewPartFiles(out, name, limit)
	if err != nil {
		return
	}
	defer files.Close()

	stats.Splits = append(stats.Splits, files.Name())

	ch := make(chan Entry)
	go func() {
		defer close(ch)
		for e := range ch {
			part, newpart, err := files.Write([]byte(e.ToString()))
			if err != nil {
				return
			}

			if newpart {
				stats.Splits = append(stats.Splits, part)
			}
		}
	}()

	if err = Scan(in, ch); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}
