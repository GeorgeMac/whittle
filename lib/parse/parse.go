package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	// ErrExpectedOnePackage is returned when multiple packages are found for the
	// provided directory location
	ErrExpectedOnePackage = errors.New("expected to find one package")

	optTagMatcher = regexp.MustCompile(`opts(?:\:"(.+?)")?`)
)

// Package contains the packages name and the types
// relevant for generating against
type Package struct {
	Name  string
	Types map[string]Type
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
	Type       string
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
		return pkg, errors.Wrapf(ErrExpectedOnePackage, "found %d packages", len(parsed))
	}

	// loop over found packages in dir
	for name, ppkg := range parsed {
		pkg = Package{Name: name, Types: map[string]Type{}}

		for _, fi := range ppkg.Files {
			// loop over files declarations
			for _, decl := range fi.Decls {
				// if declaration is a general declaration
				if gdecl, ok := decl.(*ast.GenDecl); ok {
					// loop over general declarations specs
					for _, spec := range gdecl.Specs {
						// if the spec is a type definition
						if tspec, ok := spec.(*ast.TypeSpec); ok {
							// if spec is for a struct type
							if str, ok := tspec.Type.(*ast.StructType); ok {
								pkg.Types[tspec.Name.Name] = Type{
									Name:   tspec.Name.Name,
									Fields: fetchFields(str.Fields.List),
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

func fetchFields(external []*ast.Field) (fields []Field) {
	for _, field := range external {
		if tag := field.Tag; tag != nil {
			if parts := optTagMatcher.FindStringSubmatch(tag.Value); len(parts) > 1 {
				var (
					name   = field.Names[0].Name
					method = parts[1]
				)

				if method == "" {
					// field -> WithField
					method = fmt.Sprintf("With%s", strings.Title(name))
				}

				fields = append(fields, Field{
					Name:       name,
					Type:       typeString(field.Type),
					OptionName: method,
				})
			}
		}
	}

	return
}

func typeString(e ast.Expr) string {
	switch typ := e.(type) {
	case *ast.Ident:
		// catches type identifiers like int, string and bool
		return typ.String()
	case *ast.MapType:
		key := typeString(typ.Key)
		value := typeString(typ.Value)

		return fmt.Sprintf("map[%s]%s", key, value)
	case *ast.ArrayType:
		elemType := typeString(typ.Elt)

		return fmt.Sprintf("[]%s", elemType)
	}

	return "unknown!"
}
