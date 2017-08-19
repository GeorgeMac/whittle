package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrTypeNotFound       = errors.New("type not found")
	ErrExpectedOnePackage = errors.New("expected to find one package")
)

// Package contains the packages name and the types
// relevant for generating against
type Package struct {
	Name  string
	Types []Type
}

// Type represents a struct type with a name
// and some associated fields
type Type struct {
	Name   string
	Fields []Field
}

// Field represents a field in a struct which have
// a name and function option name to be used
type Field struct {
	Name       string
	OptionName string
}

// Parse parses a director, expecting to find a single package
// It returns a Go representation of the package with the minimal
// information required for generating options code
func Parse(dir string, types ...string) (pkg Package, err error) {
	fset := token.NewFileSet()

	parsed, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return pkg, err
	}

	if len(parsed) != 1 {
		return pkg, errors.Wrapf(ErrExpectedOnePackage, "found %d", len(parsed))
	}

	// loop over found packages in dir
	for name, ppkg := range parsed {
		pkg = Package{Name: name}

		for _, fi := range ppkg.Files {
			// loop over files declerations
			for _, decl := range fi.Decls {
				// if decleration is a general declaration
				if gdecl, ok := decl.(*ast.GenDecl); ok {
					// loop over general declarations specs
					for _, spec := range gdecl.Specs {
						// if the spec is a type definition
						if tspec, ok := spec.(*ast.TypeSpec); ok {
							// if spec is for a struct type
							if _, ok := tspec.Type.(*ast.StructType); ok {
								remainingTypes := remove(types, tspec.Name.Name)
								// if we found a type
								if len(types) > len(remainingTypes) {
									types = remainingTypes

									pkg.Types = append(pkg.Types, Type{
										Name: tspec.Name.Name,
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return
}

func remove(haystack []string, needle string) []string {
	for i, s := range haystack {
		if s == needle {
			slice := make([]string, 0, len(haystack))
			slice = append(slice, haystack[0:i]...)
			slice = append(slice, haystack[i+1:]...)
			return slice
		}
	}

	return haystack
}
