package static

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

const BLOCK_SIZE = 4096

func writeGziped(fout *os.File, fin *os.File, indent string, columns uint) (string, error) {
	var buf bytes.Buffer

	h := sha1.New()

	// compressor
	z, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer z.Close()

	// read, compress and detect Content-Type
	b := make([]byte, BLOCK_SIZE)
	for {
		var n int
		b = b[:cap(b)]
		n, err = fin.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		z.Write(b[:n])
		h.Write(b[:n])
	}
	z.Flush()
	z.Close()

	// data
	bytes := buf.Bytes()
	last := len(bytes) - 1
	for i, b := range bytes {
		var pre, post string
		var col = uint(i) % columns

		if col != columns-1 {
			post = ""
		} else if i == last {
			post = ",\n" + indent
		} else {
			post = ",\n"
		}

		if col == 0 {
			// first column
			pre = "\t" + indent
		} else {
			pre = ", "
		}

		if _, err := fmt.Fprintf(fout, "%s0x%02x%s", pre, b, post); err != nil {
			return "", err
		}
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash, nil
}
