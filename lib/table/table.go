package table

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/georgemac/whittle/lib/format"
)

// Table is a type which can generate a test and
// table definition for a package, type and func
// name
type Table struct {
	Package   string
	Type      string
	Functions []string
}

// New configures a Table type definition
func New(pkg, typ string, funcs ...string) Table {
	return Table{
		Package:   pkg,
		Type:      typ,
		Functions: funcs,
	}
}

// CaseType returns the name of the test case type
// to be generated
func (t Table) CaseType(fn string) string { return fmt.Sprintf("%s%sCase", strings.ToLower(t.Type), fn) }

// WriteTestTo writes out the test definition to the writer
func (t Table) WriteTestTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	if err := testTmpl.Execute(buf, t); err != nil {
		return 0, err
	}

	return format.To(w, t.Package, buf.Bytes())
}

// WriteDefTo writes out the table definition to the writer
func (t Table) WriteDefTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	if err := defTmpl.Execute(buf, t); err != nil {
		return 0, err
	}

	return format.To(w, t.Package, buf.Bytes())
}

var (
	testTmpl = template.Must(template.New("test").Parse(`package {{ .Package }}

import "testing"
{{ range .Functions }}

func Test{{ $.Type }}{{ . }}(t *testing.T) {
	for _, testCase := range []{{ $.CaseType . }}{
		{name: "happy path"},
	} {
		t.Run(testCase.name, testCase.Run)
	}
}
{{ end }}`))

	defTmpl = template.Must(template.New("tableDef").Parse(`package {{ .Package }} 

import "testing"
{{ range .Functions }}

type {{ $.CaseType . }} struct {
	name string
}

func (c {{ $.CaseType . }}) Run(t *testing.T) {}
{{ end }}`))
)
