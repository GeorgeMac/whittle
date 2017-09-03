package options

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/georgemac/whittle/lib/options"
	"github.com/georgemac/whittle/lib/parse"
	"github.com/pkg/errors"
)

var (
	ErrTypeNotProvided = errors.New("type must be set")
	ErrTypeNotFound    = errors.New("type not found")
	ErrUsage           = errors.New("user requested usage")
)

type Command struct {
	flags *flag.FlagSet
	typ   string
}

func Parse(args []string) (Command, error) {
	var (
		command Command
		help    bool
	)

	command.flags = flag.NewFlagSet("options", flag.ContinueOnError)
	command.discardOutput()
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

func (c Command) discardOutput() {
	c.flags.SetOutput(ioutil.Discard)
}

func (c Command) Usage() {
	defer c.discardOutput()
	c.flags.SetOutput(os.Stderr)
	fmt.Println("options <options>")
	c.flags.Usage()
}

func (c Command) Run() error {
	if c.typ == "" {
		return ErrTypeNotProvided
	}

	pkg, err := parse.Parse(".", c.typ)
	if err != nil {
		return err
	}

	structType, ok := pkg.Types[c.typ]
	if !ok {
		return ErrTypeNotFound
	}

	funcs := []options.Option{}
	for _, field := range structType.Fields {
		funcs = append(funcs, options.Option{
			Name:     field.OptionName,
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
