package fasta

import (
	"bytes"
	"encoding/binary"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	"github.com/dgraph-io/badger/v3"
	"github.com/jbrough/leucine/io"
	"github.com/jbrough/leucine/seq"
)

const BatchSize = 20000

func NewIndex(db *badger.DB) (idx *Index) {
	idx = &Index{
		db:   db,
		dict: seq.NewDict(),
	}

	return
}

type Index struct {
	db   *badger.DB
	dict *seq.Dict
}

func (idx *Index) Add(src string, out chan uint64) (err error) {
	paths, err := io.PathsFromOpt(src)
	if err != nil {
		return
	}

	for _, src := range paths {
		if err := idx.add(src, out); err != nil {
			return err
		}
	}

	return
}

func (idx *Index) Align(qseq string, out chan Alignment) (err error) {
	bmq := seq.Bitmap(qseq, idx.dict)

	err = idx.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := make([]byte, 1)
		prefix[0] = uint8(1)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			_ = k
			bm := roaring.New()
			err := item.Value(func(v []byte) error {
				if err := bm.UnmarshalBinary(v); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			bmi := roaring.And(bmq, bm)

			defs := func() (r []string, err error) {
				opts := badger.DefaultIteratorOptions
				opts.PrefetchSize = 10
				it := txn.NewIterator(opts)
				defer it.Close()

				prefix := []byte{uint8(0)}
				prefix = append(prefix, k[1:9]...)
				for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
					item := it.Item()

					err := item.Value(func(v []byte) error {
						r = append(r, string(v))
						return nil
					})
					if err != nil {
						return nil, err
					}
				}
				return
			}

			if bmi.GetCardinality() > 10 {
				s := seq.Unpack(bytes.NewBuffer(k[9:]), idx.dict)

				lm := ""
				for _, i := range bmi.ToArray() {
					a := Align(qseq, s, idx.dict.Tg(uint16(i)))
					if a.Score > 30 && a.Score < 50 {
						m := strings.TrimSpace(a.Match)
						if m == lm {
							lm = m
							continue
						} else {
							lm = m
						}
						d, err := defs()
						if err != nil {
							panic(err)
						}
						a.Subject.Defs = d
						out <- a
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func (idx *Index) add(src string, out chan uint64) (err error) {
	buf := new(bytes.Buffer)
	ch := make(chan Entry)
	wg := sync.WaitGroup{}

	type rec struct {
		k []byte
		v []byte
		b uint8
	}

	dbCh := make(chan rec, BatchSize)
	var i uint64

	go func() {
		tmp := []rec{}

		save := func(tmp []rec) {
			txn := idx.db.NewTransaction(true)
			for _, r := range tmp {
				b := make([]byte, 1)
				b[0] = r.b
				k := append(b, r.k...)
				e := badger.NewEntry(k, r.v).WithMeta(byte(r.b))
				if err := txn.SetEntry(e); err == badger.ErrTxnTooBig {
					if err := txn.Commit(); err != nil {
						panic(err)
					}
					txn = idx.db.NewTransaction(true)
					if err := txn.SetEntry(e); err != nil {
						panic(err)
					} else {
						defer wg.Done()
					}
				} else {
					defer wg.Done()
				}
			}
			if err := txn.Commit(); err != nil {
				panic(err)
			}
		}

		for r := range dbCh {
			tmp = append(tmp, r)
			wg.Add(1)

			if len(tmp) == BatchSize {
				save(tmp)
				tmp = []rec{}
			}
		}
		save(tmp)
	}()

	go func() {
		defer close(dbCh)
		for e := range ch {
			wg.Add(1)

			bm := seq.Bitmap(e.Seq, idx.dict)

			bmb, err := bm.ToBytes()
			if err != nil {
				panic(err)
			}

			b := seq.Pack(e.Seq, idx.dict, buf)
			buf.Reset()
			hb := make([]byte, 8)
			binary.LittleEndian.PutUint64(hb, e.SeqHash())

			hi := make([]byte, 8)
			binary.LittleEndian.PutUint64(hi, e.Id())

			dbCh <- rec{append(hb, hi...), []byte(e.Def), 0}
			dbCh <- rec{append(hb, b...), bmb, 1}
			out <- i

			i++
			wg.Done()
		}
	}()

	err = Scan(src, ch)
	close(ch)
	wg.Wait()

	return
}
