package options

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/georgemac/whittle/lib/options"
	"github.com/georgemac/whittle/lib/parse"
	"github.com/pkg/errors"
)

var (
	// ErrTypeNotProvided is returned when the type is empty
	ErrTypeNotProvided = errors.New("type must be set")
	// ErrTypeNotFound is returned when the type cant be found
	// in the parsed definition
	ErrTypeNotFound = errors.New("type not found")
	// ErrUsage is returned when the help command is set and the
	// usage should be printed
	ErrUsage = errors.New("user requested usage")
)

// Command is the structure representation of the options command
type Command struct {
	flags *flag.FlagSet
	typ   string
}

// Parse reads the slice of arguments and returns the executable Command
func Parse(args []string) (Command, error) {
	var (
		command Command
		help    bool
	)

	command.flags = flag.NewFlagSet("options", flag.ContinueOnError)
	command.flags.SetOutput(ioutil.Discard)
	command.flags.StringVar(&command.typ, "type", "", "type for options to be generated for")
	command.flags.BoolVar(&help, "help", false, "print usage")

	if err := command.flags.Parse(args); err != nil {
		return command, errors.Wrap(err, "options")
	}

	if help {
		return command, ErrUsage
	}

	return command, nil
}

// Usage prints the flags usage and command name to Stderr
func (c Command) Usage() string {
	defer c.flags.SetOutput(ioutil.Discard)

	buf := &bytes.Buffer{}
	c.flags.SetOutput(buf)

	fmt.Fprintln(buf, "options <options>")
	c.flags.Usage()

	return buf.String()
}

// Run executes the options command
func (c Command) Run() error {
	if c.typ == "" {
		return ErrTypeNotProvided
	}

	pkg, err := parse.Parse(".")
	if err != nil {
		return err
	}

	structType, ok := pkg.Types[c.typ]
	if !ok {
		return ErrTypeNotFound
	}

	funcs := []options.Option{}
	for _, field := range structType.Fields {
		method, ok := field.Tags["opts"]
		if !ok {
			continue
		}

		if method == "" {
			method = fmt.Sprintf("With%s", strings.Title(field.Name))
		}

		funcs = append(funcs, options.Option{
			Name:     method,
			Type:     field.Type,
			Variable: field.Name,
		})
	}

	fi, err := os.Create(fmt.Sprintf("./%s_options.go", pkg.Name))
	if err != nil {
		return errors.Wrap(err, "options")
	}

	if _, err := options.New(pkg.Name, structType.Name, funcs...).WriteTo(fi); err != nil {
		return errors.Wrap(err, "options")
	}

	return nil
}
