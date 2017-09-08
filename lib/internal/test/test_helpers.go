package test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	dir = "../internal/test/testdata"
)

// FixtureFile is a structure which contains the name
// input contents and expected output contents fixtures
type FixtureFile struct {
	Name    string
	Input   string
	Outputs map[string]string
}

// StageFixture writes the input contents of the fixture
// to "{{ name }}.go" in a temporary directory and returns
// the path to this directory
func StageFixture(t *testing.T, fi FixtureFile) string {
	dir, err := ioutil.TempDir("", fi.Name)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(dir, fi.Name), []byte(fi.Input), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

// Fixture loads the input and output files for a given
// fixture name
func Fixture(t *testing.T, input string, outputs ...string) FixtureFile {
	file := FixtureFile{
		Name:    input,
		Input:   mustRead(t, filepath.Join(dir, input)),
		Outputs: map[string]string{},
	}

	for _, outPath := range outputs {
		path := filepath.Join(dir, outPath)

		file.Outputs[outPath] = mustRead(t, path)
	}

	return file
}

func mustRead(t *testing.T, path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}
