package main

import (
	"flag"
	"log"

	"go.sancus.dev/file2go/render"
)

func main() {
	c := render.Config{}
	flag.StringVar(&c.Package, "p", "", "package name")
	flag.StringVar(&c.Output, "o", "", "output file")
	flag.StringVar(&c.Varname, "N", "", "variable name")
	flag.StringVar(&c.Template, "T", "", "template type")
	flag.Parse()

	if err := c.Render(flag.Args()); err != nil {
		log.Fatal(err)
	}
}
