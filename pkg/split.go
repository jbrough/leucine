package leucine

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SplitInfo struct {
	Sources     string       `json:"source"`
	Destination string       `json:"destination"`
	Stats       []SplitStats `json:"stats"`
	RuntimeSecs float64      `json:"runtime_secs"`
}

type SplitStats struct {
	Source      string   `json:"source"`
	Splits      []string `json:"splits"`
	RuntimeSecs float64  `json:"runtime_secs"`
}

func (s SplitStats) AsJSON() string {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func SplitFasta(in, out string, limit int) (stats SplitStats, err error) {
	ts := time.Now()

	stats.Source = in
	name := strings.TrimSuffix(filepath.Base(in), ".fasta")

	file, err := os.Open(in)
	if err != nil {
		return
	}
	defer file.Close()

	files, err := NewPartFiles(out, name)
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
