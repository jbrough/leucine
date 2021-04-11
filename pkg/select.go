package blastr

import (
	"bufio"
	"os"
	"strings"
)

func Select(path, query string, out chan [2]string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	d := true

	var desc string
	var match bool
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		t := scanner.Text()
		if d {
			if strings.Contains(t, query) {
				desc = t
				match = true
			}
		} else if match {
			out <- [2]string{desc, t}

			match = false
		}
		d = !d
	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}
