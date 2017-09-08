package table

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/georgemac/whittle/lib/format"
)

type Table struct {
	Package  string
	Type     string
	Function string
}

func New(pkg, typ, fn string) Table {
	return Table{
		Package:  pkg,
		Type:     typ,
		Function: fn,
	}
}

func (t Table) CaseType() string { return fmt.Sprintf("%s%sCase", strings.ToLower(t.Type), t.Function) }

func (t Table) WriteTestTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	if err := testTmpl.Execute(buf, t); err != nil {
		return 0, err
	}

	return format.FormatTo(w, t.Package, buf.Bytes())
}

func (t Table) WriteDefTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	if err := defTmpl.Execute(buf, t); err != nil {
		return 0, err
	}

	return format.FormatTo(w, t.Package, buf.Bytes())
}

var (
	testTmpl = template.Must(template.New("test").Parse(`package {{ .Package }}

import "testing"

func Test{{ .Type }}{{ .Function }}(t *testing.T) {
	for _, testCase := range []{{ .CaseType }}{
		{name: "happy path"},
	} {
		t.Run(testCase.name, testCase.Run)
	}
}`))

	defTmpl = template.Must(template.New("tableDef").Parse(`package {{ .Package }} 

import "testing"

type {{ .CaseType }} struct {
	name string
}

func (c {{ .CaseType }}) Run(t *testing.T) {}`))
)
