package format

import (
	"go/format"
	"io"

	"golang.org/x/tools/imports"
)

// To rewrites the source data formatted with corrected importants
// to the provider writer w
func To(w io.Writer, pkg string, data []byte) (int64, error) {
	var err error

	data, err = format.Source(data)
	if err != nil {
		return 0, err
	}

	data, err = imports.Process(pkg, data, nil)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	return int64(n), err
}
