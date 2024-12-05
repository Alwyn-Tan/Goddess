package server

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type ServerConfig struct {
	Address        string
	MaxConnections uint32
	Timeout        time.Duration
}

func Serve(config *ServerConfig, handler Handler) {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Fatalf("listen err: %v", err)
	}
	var closing atomic.Bool
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		signal := <-signalChan
		switch signal {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Printf("receive exit signal: %v", signal)
			closing.Store(true)
			listener.Close()
		}
	}()

	log.Printf("server listening at %v", listener.Addr())
	defer handler.Close()
	defer listener.Close()

	ctx, _ := context.WithCancel(context.Background())
	for {
		conn, err := listener.Accept()
		if err != nil {
			if closing.Load() {
				return
			}
			log.Printf("accept err: %v", err)
			continue
		}
		log.Printf("accepted new connection from %v", conn.RemoteAddr())
		go handler.Handle(ctx, conn)
	}

}
