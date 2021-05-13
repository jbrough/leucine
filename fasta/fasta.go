package fasta

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/twmb/murmur3"
)

func NewEntry(l string) (e Entry, err error) {
	i := strings.Index(l, "\n")
	if i == -1 {
		err = errors.New("Fasta Parse Error: No NL deliminator: " + l)
		return
	}
	e = Entry{
		Def: l[:i],
		Seq: l[i+1:],
	}

	return
}

type Source struct {
	Bitmap     *roaring.Bitmap
	ByteOffset int64
	Bytes      int64
}

type Entry struct {
	Def string
	Seq string
	Src Source
}

func (e Entry) Id() uint64 {
	hasher := murmur3.New128()
	hasher.Write([]byte(e.Def))
	v1, v2 := hasher.Sum128()
	_ = v2
	return v1
}

func (e Entry) SeqHash() uint64 {
	hasher := murmur3.New128()
	hasher.Write([]byte(e.Seq))
	v1, v2 := hasher.Sum128()
	_ = v2
	return v1
}

func (e Entry) ToString() string {
	return e.Def + "\n" + e.Seq
}

func Scan(src string, out chan<- Entry) (err error) {
	var scanner *bufio.Scanner
	if strings.Index(src, "http") == 0 {
		client := new(http.Client)
		req, err := http.NewRequest("GET", src, nil)
		if err != nil {
			return err
		}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var reader io.ReadCloser

		if strings.HasSuffix(src, ".gz") {
			reader, err = gzip.NewReader(res.Body)
			if err != nil {
				panic(err)
			}
		}

		scanner = bufio.NewScanner(reader)

	} else {
		f, err := os.Open(strings.TrimPrefix(src, "file://"))
		if err != nil {
			return err
		}
		defer f.Close()

		var reader io.ReadCloser
		reader = f
		if strings.HasSuffix(src, ".gz") {
			reader, err = gzip.NewReader(f)
			if err != nil {
				panic(err)
			}
		}

		scanner = bufio.NewScanner(reader)
	}

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	return Parse(scanner, out)
}

func ScanFastaLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && data[0] == '>' {
		return 1, data, nil
	}

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	i := bytes.IndexByte(data[1:], '>')
	if i != -1 {
		if data[i] == '\n' {
			return i + 1, data[0:i], nil
		} else {
			offset := i
			for i != -1 {
				offset += 3
				i = bytes.IndexByte(data[offset:], '>')
				if data[offset+i-1] == '\n' {
					return offset + i + 1, data[0 : offset+i], nil
				}
			}
		}
	}

	// Request more data.
	return 0, nil, nil
}

func Parse(scanner *bufio.Scanner, entries chan<- Entry) (err error) {
	scanner.Split(ScanFastaLines)

	var bytesRead int64

	for scanner.Scan() {
		b := scanner.Bytes()
		l := len(b) + 1
		i := bytes.IndexByte(b, '\n')
		//for ii, a := range b[i+1:] {
		//	if a == '\n' {
		//		b[i+1+ii] = byte(0)
		//	}
		//}
		//s := string(b)
		s := string(b[:i+1]) + strings.Replace(string(b[i+1:]), "\n", "", -1)
		entries <- Entry{
			Def: s[:i],
			Seq: s[i+1:],
			Src: Source{
				ByteOffset: bytesRead,
				Bytes:      int64(l),
			},
		}
		bytesRead += int64(l)
	}

	return scanner.Err()
}
