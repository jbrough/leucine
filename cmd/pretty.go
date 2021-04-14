// Pretty print JSON output. The 'jq' package is also good for this.
// code from https://stackoverflow.com/a/53124485

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("cat file.json | prettyjson")
		return
	}

	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for {
		data := map[string]interface{}{}
		if err := dec.Decode(&data); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("decode error %v", err)
		}
		if err := enc.Encode(data); err != nil {
			log.Fatalf("encod error %v", err)
		}
	}
}
