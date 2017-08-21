package options

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/whittlingo/fun/lib/internal/test"
)

type TestCase struct {
	// input: name of the test case from which
	// the fixtures to be tested will be
	// derived and the package name
	name string
	// input: type names
	typ string
	// input: functional options to render
	funcs []Option
	// output: an expected error (if their is one expected)
	expectedError error
}

func TestOptions(t *testing.T) {
	for _, test := range []TestCase{
		{"important", "Important", []Option{
			{
				Name:     "WithField",
				Type:     "string",
				Variable: "field",
			},
			{
				Name:     "WithAttribute",
				Type:     "int",
				Variable: "attribute",
			},
			{
				Name:     "WithThings",
				Type:     "map[string]string",
				Variable: "mapOfThings",
			},
		}, nil},
	} {
		t.Run(test.name, test.Run)
	}
}

func (tc TestCase) Run(t *testing.T) {
	var (
		// read expected output fixture file
		fi = test.Fixture(t, tc.name)
		// construct options type to test
		options = New(tc.name, tc.typ, tc.funcs...)
		// output buffer to populate
		output = &bytes.Buffer{}
	)

	// write options to output buffer
	_, err := options.WriteTo(output)
	if tc.expectedError != nil {
		// if an error is expected check it
		require.Equal(t, tc.expectedError, err)
		assert.Empty(t, output.String())
	} else {
		// require the error to be nil
		require.Nil(t, err)
		assert.Equal(t, fi.Output, output.String())
	}
}

func mustRead(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func fixturePath(path string) string {
	return filepath.Join("testdata", path)
}
