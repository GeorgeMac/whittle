package options

import (
	"bytes"
	"testing"

	"github.com/georgemac/whittle/lib/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	// input: name of the test case from which
	// the fixtures to be tested will be
	// derived and the package name
	name   string
	output string
	// input: type names
	typ string
	// input: functional options to render
	funcs []Option
	// output: an expected error (if their is one expected)
	expectedError error
}

func TestOptions(t *testing.T) {
	for _, test := range []TestCase{
		{"important", "important_options.go", "Important", []Option{
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
			{
				Name:     "WithPointerToThing",
				Type:     "*string",
				Variable: "pointerToThing",
			},
			{
				Name:     "WithPointerToStruct",
				Type:     "*os.File",
				Variable: "pointerToStruct",
			},
		}, nil},
	} {
		t.Run(test.name, test.Run)
	}
}

func (tc TestCase) Run(t *testing.T) {
	var (
		// read expected output fixture file
		fi             = test.Fixture(t, tc.name+".go", tc.output)
		expectedOutput = fi.Outputs[tc.output]
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
		generated := output.String()
		assert.Equal(t, expectedOutput, generated)
	}
}
