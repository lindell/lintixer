package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lindell/lintixer/fixer"
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [path]\n", os.Args[0])

	flag.PrintDefaults()
}

func main() {
	verbose := flag.Bool("verbose", false, "print logging statements")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		return
	}
	path := args[0]

	options := []fixer.Option{
		fixer.WithNodeFixers(fixer.NonCapitalError),
	}
	if *verbose {
		options = append(options, fixer.WithLogger(&logger{}))
	}
	fixer := fixer.New(options...)

	err := fixer.Fix(path)
	if err != nil {
		fmt.Println(err)
	}
}

type logger struct{}

func (l *logger) Info(str string) {
	fmt.Println(str)
}
