package readutil

import (
	"bytes"
	"io"
)

type packagesWrapper struct {
	pkgs io.Closer
	wrap io.Reader
}

// newOptionsWrapper translates a set of packages into the
// format of the indexer
//
// set of packages:
//
//	{
//		"pkg1": { ... },
//		"pkg2": { ... }
//	}
//
// what indexer expects:
//
//	{
//		"packages": {
//		  "pkg1": { ... },
//		  "pkg1": { ... }
//		}
//	}
func PackagesWrapper(rd io.ReadCloser) *packagesWrapper {
	mrd := io.MultiReader(
		bytes.NewBufferString(`{"packages":`),
		rd,
		bytes.NewBufferString(`}`),
	)
	return &packagesWrapper{
		pkgs: rd,
		wrap: mrd,
	}
}

func (w *packagesWrapper) Read(p []byte) (n int, err error) {
	return w.wrap.Read(p)
}

func (w *packagesWrapper) Close() error {
	return w.pkgs.Close()
}
