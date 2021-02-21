package config

import (
	"fmt"
)

type Render struct {
	Package  string
	Output   string
	Template string
}

func (c *Render) Validate() error {

	switch c.Template {
	case "static", "none", "":
		c.Template = "static"
	default:
		return fmt.Errorf("Invalid Template mode %q", c.Template)
	}

	if len(c.Package) == 0 {
		return fmt.Errorf("Package name missing")
	}

	return nil
}
