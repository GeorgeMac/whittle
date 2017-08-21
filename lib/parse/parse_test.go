package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/whittlingo/fun/lib/internal/test"
)

type TestCase struct {
	name string
	typ  string
	pkg  Package
	err  error
}

func TestParse(t *testing.T) {
	for _, tc := range []TestCase{
		{
			name: "important",
			typ:  "Important",
			pkg: Package{
				Name: "important",
				Types: map[string]Type{
					"Important": {
						Name: "Important",
						Fields: []Field{
							{
								Name:       "field",
								Type:       "string",
								OptionName: "WithField",
							},
							{
								Name:       "attribute",
								Type:       "int",
								OptionName: "WithAttribute",
							},
							{
								Name:       "mapOfThings",
								Type:       "map[string]string",
								OptionName: "WithThings",
							},
						},
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
