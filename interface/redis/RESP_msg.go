package redis

type RESPMsg interface {
	ToBytes() []byte
}
