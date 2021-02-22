package file

import (
	"io"
	"net/http"
	"os"
	"strings"
)

const BLOCK_SIZE = 4096

type Blob struct {
	Body        []byte
	ContentType string
	Encoding    string
}

func (o *Blob) Load(fname string, e Encoder) (err error) {
	var fin *os.File
	var fi os.FileInfo
	var content_type string

	// open
	if fin, err = os.Open(fname); err != nil {
		return
	}
	defer fin.Close()

	// content_type
	if fi, err = fin.Stat(); err != nil {
		return err
	} else if fi.Size() == 0 {
		content_type = "application/x-empty"
	} else if strings.HasSuffix(fname, ".js") {
		content_type = "application/javascript; charset=utf-8"
	} else if strings.HasSuffix(fname, ".css") {
		content_type = "text/css; charset=utf-8"
	} else {
		b := make([]byte, 512)

		if n, err := fin.Read(b); err != nil {
			return err
		} else {
			fin.Seek(0, 0)
			content_type = http.DetectContentType(b[:n])
		}
	}

	// New
	o.ContentType = content_type

	if err = e.Reset(); err != nil {
		return err
	}

	// read and encode
	b := make([]byte, BLOCK_SIZE)
	for {
		var n int
		b = b[:cap(b)]
		n, err = fin.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		e.Write(b[:n])
	}
	e.Close()

	o.Body = e.Bytes()
	o.Encoding = e.Encoding()

	return nil
}
