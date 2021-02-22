package file

type Encoder interface {
	Encoding() string
	Bytes() []byte

	Write([]byte) (int, error)

	Reset() error
	Close() error
}
