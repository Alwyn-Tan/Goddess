package msg

type Msg interface {
	ToBytes() []byte
}
