package io

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TODO: refactor, one of those methods you split off too soon and then special case
func PathsFromOpt(opt string) (fastas []string, err error) {
	if strings.HasSuffix(opt, ".seq") || strings.HasSuffix(opt, ".fa") || strings.HasSuffix(opt, ".fasta") || strings.HasSuffix(opt, ".faa") {
		fastas = append(fastas, opt)
	} else {
		files, err := ioutil.ReadDir(opt)
		if err != nil {
			return fastas, err
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".seq") || strings.HasSuffix(f.Name(), ".fa") || strings.HasSuffix(f.Name(), ".fasta") {
				fastas = append(fastas, filepath.Join(opt, f.Name()))
			}
		}
	}

	return
}
