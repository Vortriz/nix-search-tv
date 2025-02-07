package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dgraph-io/badger/v4"
)

type Badger struct {
	badger *badger.DB
}

func NewBadger(dir string) (*Badger, error) {
	opts := badger.
		DefaultOptions(dir).
		WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open badger: %w", err)
	}

	return &Badger{
		badger: db,
	}, nil
}

func (indexer *Badger) Index(data io.Reader, indexedKeys io.Writer) error {
	pkgs := struct {
		Packages map[string]json.RawMessage `json:"packages"`
	}{}
	err := json.NewDecoder(data).Decode(&pkgs)
	if err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	batch := indexer.badger.NewWriteBatch()
	for name, pkg := range pkgs.Packages {
		nameb := []byte(name)
		err = batch.Set(nameb, injectKey(name, pkg))
		if err != nil {
			return err
		}
		indexedKeys.Write(append(nameb, []byte("\n")...))
	}

	return batch.Flush()
}

func (bdg *Badger) Load(pkgName string) (json.RawMessage, error) {
	pkg := []byte{}

	err := bdg.badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(pkgName)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()

			err := item.Value(func(v []byte) error {
				if bytes.Equal(k, prefix) {
					pkg = v
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("load value: %w", err)
			}
		}
		return nil
	})

	return pkg, err
}

func (bdg *Badger) Close() error {
	return bdg.badger.Close()
}

func injectKey(key string, pkg json.RawMessage) json.RawMessage {
	return append([]byte(`{"_key":"`+key+`",`), pkg[1:]...)
}
