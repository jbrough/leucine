package io

import (
	"fmt"
	"os"
	"path/filepath"
)

func NewPartFiles(folder, name string, limit int) (pf *PartFiles, err error) {
	pf = &PartFiles{
		Folder: folder, BaseName: name, Limit: limit}
	err = pf.Cycle()
	return
}

type PartFiles struct {
	files    []*os.File
	Folder   string
	BaseName string
	Limit    int
	entries  int
	part     int
}

func (pf *PartFiles) Name() string {
	return fmt.Sprintf("%s.%d.fa", pf.BaseName, pf.part)
}

// TODO: Buffer NewLine() and Write calls to file
func (pf *PartFiles) NewLine() (err error) {
	_, err = pf.files[len(pf.files)-1].Write([]byte("\n"))
	return
}

func (pf *PartFiles) Write(b []byte) (part string, newpart bool, err error) {
	_, err = pf.files[len(pf.files)-1].Write(b)
	if err != nil {
		return
	}

	pf.entries++
	if pf.entries%pf.Limit == 0 {
		if err = pf.Cycle(); err != nil {
			return
		}

		newpart = true
	}

	part = pf.Name()
	return
}

func (pf *PartFiles) Close() (err error) {
	if len(pf.files) > 0 {
		if err = pf.files[len(pf.files)-1].Close(); err != nil {
			return
		}
	}

	return
}

func (pf *PartFiles) Cycle() (err error) {
	if err = pf.Close(); err != nil {
		return
	}
	pf.part++
	file_path := filepath.Join(pf.Folder, pf.Name())
	f, err := os.Create(file_path)
	if err != nil {
		return err
	}
	pf.files = append(pf.files, f)

	return
}
