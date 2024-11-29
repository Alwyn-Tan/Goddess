package redis

import "net"

type Connection struct {
	conn net.Conn
}

func StartNewConnection(conn net.Conn) *Connection {
	return nil
}

func (c *Connection) Close() error {
	return nil
}
