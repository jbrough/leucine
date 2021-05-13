package io

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func PathsFromOpt(opt string) (paths []string, err error) {
	allowed_exts := []string{
		".seq", ".fa", ".fasta", ".faa", ".ent", ".sasa", ".gz",
	}

	hasSuffix := func(s string) bool {
		for _, ext := range allowed_exts {
			if strings.HasSuffix(s, ext) {
				return true
			}
		}
		return false
	}

	if hasSuffix(opt) {
		paths = append(paths, opt)
	} else {
		files, err := ioutil.ReadDir(opt)
		if err != nil {
			return paths, err
		}

		for _, f := range files {
			if hasSuffix(f.Name()) {
				paths = append(paths, filepath.Join(opt, f.Name()))
			}
		}
	}

	return
}
