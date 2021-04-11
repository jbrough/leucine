package blastr

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createPartFile(folder, name string, part int) (f *os.File, err error) {
	file_name := fmt.Sprintf("%s.%d.fa", name, part)
	file_path := filepath.Join(folder, file_name)

	return os.Create(file_path)
}

func SplitFasta(in, out string, limit int) (err error) {
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

	return
}
