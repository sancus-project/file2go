package static

import (
	"fmt"
	"net/http"
	"os"
)

func WriteEmpty(fout *os.File, content_type string) (err error) {
	if content_type == "" {
		content_type = "application/x-empty"
	}

	_, err = fmt.Fprintf(fout, `static.Content{
	ContentType: %q,
	Body:        []byte{},
}`, content_type)

	return
}

func WriteContent(fout *os.File, fin *os.File, content_type string) error {
	var err error
	var sha1 string

	if fi, err := fin.Stat(); err != nil {
		return err
	} else if fi.Size() == 0 {
		return WriteEmpty(fout, content_type)
	}

	if content_type == "" {
		// content_type
		b := make([]byte, 512)
		n, err := fin.Read(b)
		if err != nil {
			return err
		}
		fin.Seek(0, 0)
		content_type = http.DetectContentType(b[:n])
	}

	// header
	if _, err = fmt.Fprintf(fout, "static.Content{\n\tContentType: %q,\n\tBody: []byte{\n", content_type); err != nil {
		return err
	}

	// data
	if sha1, err = writeGziped(fout, fin, "\t", 8); err != nil {
		return err
	}

	// footer
	if _, err = fmt.Fprintf(fout, "},\n\tSha1sum: %q,\n}", sha1); err != nil {
		return err
	}

	return nil
}
