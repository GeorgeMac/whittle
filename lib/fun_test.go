package fun

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	// input: name of the test case from which
	// the fixtures to be tested will be
	// derived and the package name
	name string
	// input: type names
	typ string
	// output: an expected error (if their is one expected)
	expectedError error
}

func test(name, typ, vrb string, err error) TestCase {
	return TestCase{
		name:          name,
		typ:           typ,
		expectedError: err,
	}
}

func TestOptions(t *testing.T) {
	for _, test := range []TestCase{
		test("important", "Important", "i", nil),
	} {
		t.Run(test.name, test.Run)
	}
}

func (tc TestCase) Run(t *testing.T) {
	var (
		// e.g. important_options.go
		outputFixtureName = tc.name + "_options.go"
		// read expected output fixture file
		expectedOutput = mustRead(fixturePath(outputFixtureName))
		// construct options type to test
		options = New(tc.name, tc.typ)
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
		assert.Equal(t, expectedOutput, output.String())
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
