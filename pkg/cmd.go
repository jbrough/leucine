package leucine

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func FastaPathsFromOpt(opt string) (fastas []string, err error) {
	if strings.HasSuffix(opt, ".fa") || strings.HasSuffix(opt, ".fasta") || strings.HasSuffix(opt, ".faa") {
		fastas = append(fastas, opt)
	} else {
		files, err := ioutil.ReadDir(opt)
		if err != nil {
			return fastas, err
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".fa") || strings.HasSuffix(f.Name(), ".fasta") {
				fastas = append(fastas, filepath.Join(opt, f.Name()))
			}
		}
	}

	return
}
