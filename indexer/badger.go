package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/3timeslazy/nix-search-tv/indexer/x/jsonstream"
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

// Package defines fields set by the indexer during
// indexing
type Package struct {
	Name string `json:"_key"`
}

// Indexable represents the internal structure of the data
// that the indexer expects.
type Indexable struct {
	Packages map[string]json.RawMessage `json:"packages"`
}

func (indexer *Badger) Index(data io.Reader, indexedKeys io.Writer) error {
	// Delete previous index. If we do not do that
	// and just re-assign the keys below, the index
	// will be updated, however its size will increase drastically.
	// Then, to keep the index size small, we'll need to deal with
	// badger's garbade colletion. So, it's just easier to drop everything
	err := indexer.badger.DropAll()
	if err != nil {
		return fmt.Errorf("drop all: %w", err)
	}

	batch := indexer.badger.NewWriteBatch()

	err = jsonstream.ParsePackages(data, func(name string, content []byte) error {
		nameb := []byte(name)

		err := batch.Set(nameb, injectKey(name, content))
		if err != nil {
			return fmt.Errorf("set %s: %w", name, err)
		}

		indexedKeys.Write(append(nameb, []byte("\n")...))

		return nil
	})
	if err != nil {
		return fmt.Errorf("handle packages: %w", err)
	}

	return batch.Flush()
}

func (bdg *Badger) Load(pkgName string) (json.RawMessage, error) {
	pkg := []byte{}

	err := bdg.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(pkgName))
		if err != nil {
			return err
		}

		pkg, err = item.ValueCopy(pkg)
		if err != nil {
			return err
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
	return append([]byte(`{"_key":`+strconv.Quote(key)+`,`), pkg[1:]...)
}
