package static

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type Blob []byte

func (b Blob) NewReader() (*gzip.Reader, error) {
	// bytes.NewReader() takes ownership of the given buffer
	// but what does this actually mean?
	return gzip.NewReader(bytes.NewReader(b[:]))

}

func (b Blob) RawBytes() []byte {
	return b[:]
}

func (b Blob) Bytes() ([]byte, error) {
	fz, err := b.NewReader()
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (b Blob) String() string {
	s, err := b.Bytes()
	if err != nil {
		panic(err)
	}
	return string(s)
}

func (b Blob) Encoding() string {
	return "gzip"
}
