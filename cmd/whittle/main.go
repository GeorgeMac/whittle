package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/georgemac/whittle/cmd/whittle/options"
	"github.com/pkg/errors"
)

func printUsage() {
	fmt.Println(`whittle [cmd] <options>`)
	fmt.Println(`commands:`)
	fmt.Println("\toptions - generate functional options for a type")
}

func help() {
	printUsage()
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		help()
	}

	command := os.Args[1]
	switch command {
	case "options":
		command, err := options.Parse(os.Args[2:])
		if err != nil {
			if cause := errors.Cause(err); cause == options.ErrUsage || cause == flag.ErrHelp {
				fmt.Print("whittle ")
				command.Usage()
			} else {
				fmt.Println(err)
			}

			os.Exit(1)
		}

		if err := command.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		help()
	}
}
