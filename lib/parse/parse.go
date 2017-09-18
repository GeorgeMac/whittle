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

	tagMatcher = regexp.MustCompile(`([a-z]+)(?:\:"(.+?)")?`)
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
	Funcs  []Func
}

// Field represents a field in a struct which have
// a name and function option name to be used
type Field struct {
	Name string
	Type string
	Tags map[string]string
}

// Func represntes function definition
type Func struct {
	Name string
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
			for _, gdecl := range fi.Decls {
				switch decl := gdecl.(type) {
				// if declaration is a general declaration
				case *ast.GenDecl:
					// loop over general declarations specs
					for _, spec := range decl.Specs {
						// if the spec is a type definition
						if tspec, ok := spec.(*ast.TypeSpec); ok {
							switch typ := tspec.Type.(type) {
							case *ast.StructType:
								var (
									name   = tspec.Name.Name
									fields = fetchFields(typ.Fields.List)
								)

								if ptyp, ok := pkg.Types[name]; ok {
									ptyp.Fields = fields
									pkg.Types[name] = ptyp
									continue
								}

								pkg.Types[name] = Type{
									Name:   name,
									Fields: fields,
								}
							}
						}
					}
				case *ast.FuncDecl:
					for _, recv := range decl.Recv.List {
						var (
							typ = typeString(recv.Type)
							fn  = Func{Name: decl.Name.Name}
						)

						if len(typ) < 1 {
							continue
						}

						if typ[0] == '*' {
							typ = typ[1:]
						}

						if ptyp, ok := pkg.Types[typ]; ok {
							ptyp.Funcs = append(ptyp.Funcs, fn)
							pkg.Types[typ] = ptyp
						} else {
							pkg.Types[typ] = Type{
								Name:  typ,
								Funcs: []Func{fn},
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
		if tags := field.Tag; tags != nil {
			for _, tag := range strings.Split(tags.Value, ",") {
				if parts := tagMatcher.FindStringSubmatch(tag); len(parts) > 1 {
					var (
						name   = field.Names[0].Name
						tag    = parts[1]
						method = parts[2]
					)

					fields = append(fields, Field{
						Name: name,
						Type: typeString(field.Type),
						Tags: map[string]string{tag: method},
					})
				}
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
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", typeString(typ.X))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", typeString(typ.X), typ.Sel.String())
	}

	return "unknown!"
}
