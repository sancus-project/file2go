package file

import (
	"bytes"
	"compress/gzip"
	"fmt"
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
	} else if strings.HasSuffix(fname, ".svg") {
		content_type = "image/svg+xml; charset=utf-8"
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

func (o *Blob) NewEncodedReader() io.Reader {
	return bytes.NewReader(o.Body)
}

func (o *Blob) NewReader2(encoding string) (r io.Reader, err error) {
	r = o.NewEncodedReader()

	switch encoding {
	case "":
		// as-is
	case "gzip":
		r, err = gzip.NewReader(r)
	default:
		r = nil
		err = fmt.Errorf("file2go: unknown encoding: %q", encoding)
	}

	return
}

func (o *Blob) NewReader() (r io.Reader, err error) {
	return o.NewReader2(o.Encoding)
}

func (o *Blob) RenderBytes(fout *os.File, indent string, columns uint) error {
	var err error
	var single bool

	l := len(o.Body)
	last := l - 1

	_, err = fout.WriteString("[]byte{")
	if err != nil {
		return err
	}

	if uint(l) > columns {
		single = false
	} else {
		single = true
	}

	for i, b := range o.Body {
		var pre, post string
		var col = uint(i) % columns

		if col == 0 {
			if single {
				pre = ""
			} else {
				pre = "\n" + indent
			}
		} else {
			pre = " "
		}

		if i == last {
			post = ""
		} else {
			post = ","
		}

		_, err = fmt.Fprintf(fout, "%s0x%02x%s", pre, b, post)
		if err != nil {
			return err
		}
	}

	_, err = fout.WriteString("}")
	return err
}
