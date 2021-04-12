package leucine

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type SelectInfo struct {
	Source      string        `json:"source"`
	Destination string        `json:"destination"`
	Query       string        `json:"query"`
	Stats       []SelectStats `json:"stats"`
	RuntimeSecs float64       `json:"runtime_secs"`
}

type SelectStats struct {
	FastaFile   string  `json:"fasta_file"`
	Searched    int     `json:"searched"`
	Selected    int     `json:"selected"`
	RuntimeSecs float64 `json:"runtime_secs"`
}

func (s SelectStats) Add(sa SelectStats) SelectStats {
	return SelectStats{
		Searched: s.Searched + sa.Searched,
		Selected: s.Selected + sa.Selected,
	}
}

func Select(path, query string, out chan [2]string) (stats SelectStats, err error) {
	ts := time.Now()

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	d := true

	var desc string
	var match bool
	scanner := bufio.NewScanner(file)

	stats = SelectStats{}
	stats.FastaFile = path

	query = strings.ToLower(query)

	for scanner.Scan() {
		t := scanner.Text()
		if d {
			stats.Searched++
			if strings.Contains(strings.ToLower(t), query) {
				desc = t
				match = true
			}
		} else if match {
			out <- [2]string{desc, t}
			stats.Selected++
			match = false
		}
		d = !d
	}
	if err = scanner.Err(); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}
