package main

import (
	"flag"
	"log"

	"github.com/amery/file2go/config"
	"github.com/amery/file2go/render"
)

func main() {
	c := config.Render{}
	flag.StringVar(&c.Package, "p", "", "package name")
	flag.StringVar(&c.Output, "o", "", "output file")
	flag.StringVar(&c.Template, "T", "", "template type")
	flag.Parse()

	if err := c.Validate(); err != nil {
		log.Fatal(err)
	}

	if err := render.RenderConfig(c, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
