package metrics

import "encoding/json"

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
