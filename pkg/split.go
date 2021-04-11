package blastr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SplitInfo struct {
	Sources     string     `json:"source"`
	Destination string     `json:"destination"`
	Stats       SplitStats `json:"stats"`
	RuntimeSecs float64    `json:"runtime_secs"`
}

type SplitStats struct {
	Source      string   `json:"source"`
	Splits      []string `json:"splits"`
	RuntimeSecs float64  `json:"runtime_secs"`
}

func (s SplitStats) Add(sa SplitStats) SplitStats {
	stats := SplitStats{
		Source: s.Source + sa.Source}

	for _, sp := range s.Splits {
		stats.Splits = append(stats.Splits, sp)
	}
	for _, sp := range sa.Splits {
		stats.Splits = append(stats.Splits, sp)
	}

	return stats
}

func (s SplitStats) AsJSON() string {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func createPartFile(folder, name string, part int) (f *os.File, err error) {
	file_name := fmt.Sprintf("%s.%d.fa", name, part)
	file_path := filepath.Join(folder, file_name)
	return os.Create(file_path)
}

func SplitFasta(in, out string, limit int) (stats SplitStats, err error) {
	ts := time.Now()

	stats.Source = in
	name := strings.TrimSuffix(filepath.Base(in), ".fasta")

	part := 1

	file, err := os.Open(in)
	if err != nil {
		return
	}
	defer file.Close()

	f, err := createPartFile(out, name, part)
	if err != nil {
		return
	}
	defer f.Close()
	stats.Splits = append(stats.Splits, f.Name())

	var c bool
	scanner := bufio.NewScanner(file)
	var desc string
	var seq string
	var count int
	for scanner.Scan() {
		t := scanner.Text()
		if c {
			if strings.Contains(t, ">") {
				count++
				c = false
				line := desc + "\n" + seq + "\n"
				if _, err = f.WriteString(line); err != nil {
					return
				}
				desc = ""
				seq = ""
				if count%limit == 0 {
					f.Close()
					part++
					f, err = createPartFile(out, name, part)
					if err != nil {
						return
					}
					stats.Splits = append(stats.Splits, f.Name())
				}
			} else {
				seq = seq + t
			}
		}
		if strings.Contains(t, ">") {
			desc = t
			c = true
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	es := time.Now().Sub(ts).Seconds()
	stats.RuntimeSecs = es

	return
}
