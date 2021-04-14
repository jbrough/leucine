package metrics

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
