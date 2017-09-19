package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/georgemac/whittle/cmd/whittle/options"
	"github.com/georgemac/whittle/cmd/whittle/table"
	"github.com/pkg/errors"
)

func printUsage() {
	fmt.Println(`whittle [cmd] <flags>`)
	fmt.Println(`commands:`)
	fmt.Println("\toptions - generate functional options for a type")
}

func help() {
	printUsage()
	os.Exit(1)
}

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

var (
	optionsCommand ParserFunc = func(args []string) (Command, error) { return options.Parse(args) }
	tableCommand   ParserFunc = func(args []string) (Command, error) { return table.Parse(args) }
	commands                  = map[string]Parser{
		"options": optionsCommand,
		"table":   tableCommand,
	}
)

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
			if cause := errors.Cause(err); cause == options.ErrUsage || cause == table.ErrUsage || cause == flag.ErrHelp {
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
