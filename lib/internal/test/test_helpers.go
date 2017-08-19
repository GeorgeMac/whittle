package test

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	dir = "../internal/test/testdata"
)

type FixtureFile struct {
	Name          string
	Input, Output string
}

func StageFixture(t *testing.T, fi FixtureFile) string {
	dir, err := ioutil.TempDir("", fi.Name)
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filepath.Join(dir, fi.Name+".go"), []byte(fi.Input), 0644); err != nil {
		t.Fatal(err)
	}

	return dir
}

func Fixture(t *testing.T, name string) FixtureFile {
	inpath := filepath.Join(dir, name+".go")
	output := filepath.Join(dir, name+"_options.go")

	return FixtureFile{
		Name:   name,
		Input:  mustRead(t, inpath),
		Output: mustRead(t, output),
	}
}

func mustRead(t *testing.T, path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return string(data)
}
