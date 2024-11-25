package tcp

import (
	"context"
	"net"
)

type handle func(ctx context.Context, conn net.Conn)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
