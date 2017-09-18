package parse

import (
	"testing"

	"github.com/georgemac/whittle/lib/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	// inputs
	name string
	// expecations
	typ string
	pkg Package
	err error
}

func TestParse(t *testing.T) {
	for _, tc := range []TestCase{
		{
			name: "important.go",
			typ:  "Important",
			pkg: Package{
				Name: "important",
				Types: map[string]Type{
					"Important": {
						Name: "Important",
						Fields: []Field{
							{
								Name: "field",
								Type: "string",
								Tags: map[string]string{"opts": ""},
							},
							{
								Name: "attribute",
								Type: "int",
								Tags: map[string]string{"opts": ""},
							},
							{
								Name: "mapOfThings",
								Type: "map[string]string",
								Tags: map[string]string{"opts": "WithThings"},
							},
							{
								Name: "pointerToThing",
								Type: "*string",
								Tags: map[string]string{"opts": ""},
							},
							{
								Name: "pointerToStruct",
								Type: "*os.File",
								Tags: map[string]string{"opts": ""},
							},
						},
					},
				},
			},
		},
		{
			name: "processor.go",
			typ:  "Processor",
			pkg: Package{
				Name: "processor",
				Types: map[string]Type{
					"Processor": {
						Name:  "Processor",
						Funcs: []Func{{"Run"}},
					},
				},
			},
		},
	} {
		t.Run(tc.name, tc.Run)
	}
}

func (tc *TestCase) Run(t *testing.T) {
	var (
		// get the fixture for the test case name
		fi = test.Fixture(t, tc.name)
		// write fixtures input file to testing directory
		dir = test.StageFixture(t, fi)
	)

	// run parse with staged directory and test case type
	pkg, err := Parse(dir, tc.typ)

	// ensure response is as expected
	require.Equal(t, tc.err, err)
	assert.Equal(t, tc.pkg, pkg)
}
