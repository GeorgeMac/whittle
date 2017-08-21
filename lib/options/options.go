package options

import (
	"bytes"
	"go/format"
	"io"
	"strings"
	"text/template"
)

// Options is a type which can be written to.
// It writes out Go code representation of functional options
// for the provided package, type and variable name
type Options struct {
	Package   string
	Type      string
	Functions []Option
}

// Option is a representation of the functional option to be rendered
type Option struct {
	Name     string
	Type     string
	Variable string
}

// New returns a new Options type
func New(pkg, typ string, opts ...Option) *Options {
	return &Options{
		Package:   pkg,
		Type:      typ,
		Functions: opts,
	}
}

// Var returns the variable name which will be used
// in the generated source code
func (o Options) Var() string {
	if len(o.Type) < 1 {
		return ""
	}

	return strings.ToLower(o.Type)[:1]
}

// WriteTo serializes the options type to the writer
func (o *Options) WriteTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	if err := optionsTmpl.Execute(buf, o); err != nil {
		return 0, err
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)
	return int64(n), err
}

var (
	optionsTmpl = template.Must(template.New("options").Parse(`package {{ .Package }}

type Option func(*{{ .Type }})

type Options []Option

func (o Options) Apply({{ .Var }} *{{ .Type }}) {
    for _, opt := range o {
        opt({{ .Var }})
    }
}

{{ range .Functions }}
func {{ .Name }}({{ .Variable }} {{ .Type }}) Option {
    return func({{ $.Var }} *{{ $.Type }}) {
        {{ $.Var }}.{{ .Variable }} = {{ .Variable }}
    }
}
{{ end }}
`))
)
