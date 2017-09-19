package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/georgemac/whittle/cmd/whittle/options"
	"github.com/georgemac/whittle/cmd/whittle/table"
	"github.com/pkg/errors"
)

var (
	commands = map[string]Parser{
		"options": ParserFunc(func(args []string) (Command, error) { return options.Parse(args) }),
		"table":   ParserFunc(func(args []string) (Command, error) { return table.Parse(args) }),
	}

	usageErrors = ErrorCauses{flag.ErrHelp, options.ErrUsage, table.ErrUsage}
)

// Parser is an interface for types which parse arguments and
// returns Command types to be ran
type Parser interface {
	Parse(args []string) (Command, error)
}

// ParserFunc is a function which matches the Parser interface Parse function signature
type ParserFunc func([]string) (Command, error)

// Parse delegates the call to the receiver
func (p ParserFunc) Parse(args []string) (Command, error) { return p(args) }

// Command is something which can be Ran or a Usage string be produced upon
type Command interface {
	Run() error
	Usage() string
}

// ErrorCauses is a slice of errors
type ErrorCauses []error

// ContainsCause returns true if the cause of err is in the
// ErrorCauses slice
func (e ErrorCauses) ContainsCause(err error) bool {
	for _, cause := range e {
		if errors.Cause(err) == cause {
			return true
		}
	}

	return false
}

func printUsage() {
	fmt.Println(`whittle [cmd] <flags>`)
	fmt.Println(`commands:`)
	fmt.Println("\toptions - generate functional options for a type")
	fmt.Println("\ttable - generate table driven tests for a type")
}

func help() {
	printUsage()
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	var (
		command = os.Args[1]
		args    = os.Args[2:]
	)

	if cmd, ok := commands[command]; ok {
		command, err := cmd.Parse(args)
		if err != nil {
			if usageErrors.ContainsCause(err) {
				fmt.Print("whittle ", command.Usage())
			} else {
				fmt.Println(err)
			}

			os.Exit(1)
		}

		if err := command.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
	}

	help()
}
