package table

import (
	"bytes"
	"testing"

	"github.com/georgemac/whittle/lib/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableWriteTo(t *testing.T) {
	var (
		testDef                 = "processor_test.go"
		tableDef                = "processor_table_test.go"
		fixture                 = test.Fixture(t, "processor.go", testDef, tableDef)
		testOutput, tableOutput = &bytes.Buffer{}, &bytes.Buffer{}
	)

	output := New("processor", "Processor", "Run")

	if _, err := output.WriteTestTo(testOutput); err != nil {
		require.Nil(t, err)
	}

	assert.Equal(t, fixture.Outputs[testDef], testOutput.String())

	if _, err := output.WriteDefTo(tableOutput); err != nil {
		require.Nil(t, err)
	}

	assert.Equal(t, fixture.Outputs[tableDef], tableOutput.String())
}
