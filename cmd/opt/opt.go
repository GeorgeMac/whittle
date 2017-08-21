package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/georgemac/whittle/lib/options"
	"github.com/georgemac/whittle/lib/parse"
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
		log.Fatal(err)
	}

	if _, err := options.New(pkg.Name, structType.Name, funcs...).WriteTo(fi); err != nil {
		log.Fatal(err)
	}
}
