package main

import (
	"flag"
)

type config struct {
	Package string
	Output  string

	omitGoGenerate bool
	appendOutput   bool
}

func main() {
	c := config{}
	flag.StringVar(&c.Package, "p", "", "package name")
	flag.StringVar(&c.Output, "o", "", "output file")
	flag.BoolVar(&c.omitGoGenerate, "G", false, "omit //go:generate")
	flag.BoolVar(&c.appendOutput, "a", false, "append to existing file")
	flag.Parse()

	c.Process(flag.Args())
}
