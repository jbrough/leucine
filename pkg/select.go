package blastr

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type SelectInfo struct {
	Query       string      `json:"query"`
	Destination string      `json:"destination"`
	Stats       SelectStats `json:"stats"`
	RuntimeSecs float64     `json:"runtime_secs"`
}

type SelectStats struct {
	FastasSearched int     `json:"fastas_searched"`
	FastasSelected int     `json:"fastas_selected"`
	RuntimeSecs    float64 `json:"runtime_secs"`
}

func (s SelectStats) Add(sa SelectStats) SelectStats {
	return SelectStats{
		s.FastasSearched + sa.FastasSearched,
		s.FastasSelected + sa.FastasSelected,
		0,
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

	for scanner.Scan() {
		t := scanner.Text()
		if d {
			stats.FastasSearched++
			if strings.Contains(t, query) {
				desc = t
				match = true
			}
		} else if match {
			out <- [2]string{desc, t}
			stats.FastasSelected++
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
