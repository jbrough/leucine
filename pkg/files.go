package leucine

import (
	"fmt"
	"os"
	"path/filepath"
)

func NewPartFiles(folder, name string) (pf *PartFiles, err error) {
	pf = &PartFiles{
		Folder: folder, BaseName: name}
	err = pf.Cycle()
	return
}

type PartFiles struct {
	files    []*os.File
	Folder   string
	BaseName string
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

func (pf *PartFiles) Write(b []byte) (err error) {
	_, err = pf.files[len(pf.files)-1].Write(b)
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
