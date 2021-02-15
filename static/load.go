package static

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const BLOCK_SIZE = 4096

func NewContent(fname string) (o *Content, err error) {
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
		return nil, err
	} else if fi.Size() == 0 {
		content_type = "application/x-empty"
	} else if strings.HasSuffix(fname, ".js") {
		content_type = "application/javascript; charset=utf-8"
	} else if strings.HasSuffix(fname, ".css") {
		content_type = "text/css; charset=utf-8"
	} else {
		b := make([]byte, 512)

		if n, err := fin.Read(b); err != nil {
			return nil, err
		} else {
			fin.Seek(0, 0)
			content_type = http.DetectContentType(b[:n])
		}
	}

	// New
	o = &Content{
		ContentType: content_type,
	}

	if err = o.Load(fin); err != nil {
		return nil, err
	}

	return o, err
}

func (o *Content) Load(fin *os.File) error {
	var buf bytes.Buffer

	h := sha1.New()

	// compressor
	z, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return err
	}
	defer z.Close()

	// read, compress and hash
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
		z.Write(b[:n])
		h.Write(b[:n])
	}
	z.Flush()
	z.Close()

	// data
	o.Body = buf.Bytes()
	o.Sha1sum = fmt.Sprintf("%x", h.Sum(nil))
	return nil
}
