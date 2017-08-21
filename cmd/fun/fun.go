package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/whittlingo/fun/lib/options"
	"github.com/whittlingo/fun/lib/parse"
)

var (
	typ string
)

func main() {
	flag.StringVar(&typ, "type", "", "type for options to be generated for")
	flag.Parse()

	if typ == "" {
		log.Fatal("must provide -type")
	}

	pkg, err := parse.Parse(".", typ)
	if err != nil {
		log.Fatal(err)
	}

	structType, ok := pkg.Types[typ]
	if !ok {
		log.Fatalf("cannot find type %q", typ)
	}

	fi, err := os.Create(fmt.Sprintf("./%s_options.go", pkg.Name))
	if err != nil {
		log.Fatal(err)
	}

	if _, err := options.New(pkg.Name, structType.Name).WriteTo(fi); err != nil {
		log.Fatal(err)
	}
}
