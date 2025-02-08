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

type BadgerConfig struct {
	Dir      string
	InMemory bool
}

func NewBadger(conf BadgerConfig) (*Badger, error) {
	opts := badger.
		DefaultOptions(conf.Dir).
		WithLoggingLevel(badger.ERROR).
		WithInMemory(conf.InMemory)
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
			if !bytes.Equal(k, prefix) {
				continue
			}

			var err error
			pkg, err = item.ValueCopy(nil)
			if err != nil {
				return fmt.Errorf("copy value: %w", err)
			}

			break
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("iter failed: %w", err)
	}

	return pkg, nil
}

func (bdg *Badger) Close() error {
	return bdg.badger.Close()
}

// injectKey appends the `_key` field into the json object.
//
// This thing saves about ~2.5s on my laptop when indexing 120k nix packages
func injectKey(key string, pkg json.RawMessage) json.RawMessage {
	return append([]byte(`{"_key":"`+key+`",`), pkg[1:]...)
}
