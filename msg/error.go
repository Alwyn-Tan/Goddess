package msg

type Error interface {
	Error() string
	ToBytes() []byte
}
