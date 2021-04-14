package fasta

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/jbrough/leucine/metrics"
)

func Select(path, query string, out chan [2]string) (stats metrics.SelectStats, err error) {
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

	stats = metrics.SelectStats{}
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
