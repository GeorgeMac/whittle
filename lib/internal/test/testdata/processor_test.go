package processor

import "testing"

func TestProcessorRun(t *testing.T) {
	for _, testCase := range []processorRunCase{
		{name: "happy path"},
	} {
		t.Run(testCase.name, testCase.Run)
	}
}
