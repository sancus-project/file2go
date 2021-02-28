package template

import (
	"bytes"
	"compress/gzip"
)

type TemplateValidator func([]byte) error

type TemplateEncoder struct {
	gzipped bytes.Buffer
	plain bytes.Buffer

	z *gzip.Writer

	Validator TemplateValidator
}

func (e *TemplateEncoder) Encoding() string {
	return "gzip"
}

func (e *TemplateEncoder) Bytes() []byte {
	return e.gzipped.Bytes()
}

func (e *TemplateEncoder) Close() error {
	if e.z != nil {
		e.z.Flush()
		e.z.Close()
	}

	if e.Validator != nil {
		return e.Validator(e.plain.Bytes())
	} else {
		return nil
	}
}

func (e *TemplateEncoder) Reset() error {
	e.gzipped.Reset()
	e.plain.Reset()

	// compressor
	z, err := gzip.NewWriterLevel(&e.gzipped, gzip.BestCompression)
	if err != nil {
		return err
	}

	e.z = z
	return nil
}

func (e *TemplateEncoder) Write(b []byte) (int, error) {
	n := len(b)

	if n > 0 {
		if _, err := e.z.Write(b); err != nil {
			return 0, err
		}
	}

	return n, nil
}
