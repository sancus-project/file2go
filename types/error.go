package types

type Error interface {
	Error() string
	Status() int
}
