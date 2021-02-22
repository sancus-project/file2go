package static

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"hash"
)

// file.Encoder
type StaticRenderEncoder struct {
	buf bytes.Buffer

	h hash.Hash
	z *gzip.Writer

	sha1sum string
}

func (e *StaticRenderEncoder) Encoding() string {
	return "gzip"
}

func (e *StaticRenderEncoder) Bytes() []byte {
	return e.buf.Bytes()
}

func (e *StaticRenderEncoder) Sha1sum() string {
	return e.sha1sum
}
func (e *StaticRenderEncoder) Close() error {
	if e.z != nil {
		e.z.Flush()
		e.z.Close()
	}

	if e.h != nil {
		e.sha1sum = fmt.Sprintf("%x", e.h.Sum(nil))
	} else {
		e.sha1sum = ""
	}
	return nil
}

func (e *StaticRenderEncoder) Reset() error {
	e.buf.Reset()

	// hasher
	h := sha1.New()
	// compressor
	z, err := gzip.NewWriterLevel(&e.buf, gzip.BestCompression)
	if err != nil {
		return err
	}

	e.z = z
	e.h = h
	return nil
}

func (e *StaticRenderEncoder) Write(b []byte) (int, error) {
	n := len(b)
	if _, err := e.h.Write(b); err != nil {
		return 0, err
	}
	if _, err := e.z.Write(b); err != nil {
		return 0, err
	}
	return n, nil
}
