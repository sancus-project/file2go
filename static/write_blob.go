package static

import (
	"os"
)

func WriteBlob(fout *os.File, fin *os.File) error {
	// header
	if _, err := fout.WriteString("[]byte{\n"); err != nil {
		return err
	}

	// data
	if _, err := writeGziped(fout, fin, "", 8); err != nil {
		return err
	}

	// footer
	if _, err := fout.WriteString("}"); err != nil {
		return err
	}

	return nil
}
