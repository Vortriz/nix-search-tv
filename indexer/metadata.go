package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
)

const metadataFile = "metadata.json"
const cacheFile = "cache.txt"

func GetMetadata(dir string) (IndexMetadata, error) {
	md := IndexMetadata{}
	path, err := initFile(dir, metadataFile, []byte("{}"))
	if err != nil {
		return md, fmt.Errorf("init metadata: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return md, fmt.Errorf("read metadata: %w", err)
	}

	err = json.Unmarshal(data, &md)
	if err != nil {
		return md, fmt.Errorf("unmarshal metadata: %w", err)
	}
	return md, nil
}

func SetMetadata(dir string, md IndexMetadata) error {
	mdpath, err := initFile(dir, metadataFile, []byte("{}"))
	if err != nil {
		return fmt.Errorf("init metadata: %w", err)
	}

	data, err := json.Marshal(md)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}
	err = os.WriteFile(mdpath, data, 0666)
	if err != nil {
		return fmt.Errorf("write metadata: %w", err)
	}

	return nil
}

func CacheWriter(dir string) (io.WriteCloser, error) {
	cpath, err := initFile(dir, cacheFile, nil)
	if err != nil {
		return nil, fmt.Errorf("init cache: %w", err)
	}

	return os.OpenFile(cpath, os.O_WRONLY, 0666)
}

func CacheReader(dir string) (io.ReadCloser, error) {
	cpath, err := initFile(dir, cacheFile, nil)
	if err != nil {
		return nil, fmt.Errorf("init cache: %w", err)
	}

	return os.OpenFile(cpath, os.O_RDONLY, 0666)
}

func initFile(dir, file string, initValue []byte) (string, error) {
	mdpath := path.Join(dir, file)

	_, err := os.Stat(mdpath)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.WriteFile(mdpath, initValue, 0666)
		if err != nil {
			return "", fmt.Errorf("write empty metadata: %w", err)
		}

		return mdpath, nil
	}
	if err != nil {
		return "", fmt.Errorf("stat failed: %w", err)
	}

	return mdpath, nil
}
