package tcp

import (
	"Goddess/lib/sync/wait"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type EchoHandler struct {
	activeConnMap sync.Map
	closing       atomic.Bool
}

type Client struct {
	Conn      net.Conn
	WaitGroup wait.Wait
}

func InitEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (c *Client) Close() error {
	if !c.WaitGroup.WaitWithTimeout(10 * time.Second) {
		return fmt.Errorf("timeout while waiting")
	}
	err := c.Conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}
	return nil
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		conn.Close()
	}

	client := &Client{
		Conn: conn,
	}

	h.activeConnMap.Store(client, 1)

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("connection close")
				h.activeConnMap.Delete(conn)
			} else {
				log.Printf("error reading message: %v", err)
			}
			return
		}
		client.WaitGroup.Add(1)
		b := []byte(msg)
		err = h.Broadcast(b, conn)
		if err != nil {
			return
		}
		client.WaitGroup.Done()
	}
}

func (h *EchoHandler) Close() error {
	log.Printf("closing echo connection")
	h.closing.Store(true)
	h.activeConnMap.Range(func(key interface{}, value interface{}) bool {
		client := key.(*Client)
		client.Close()
		return true
	})
	return nil
}

func (h *EchoHandler) Broadcast(msg []byte, sender net.Conn) error {
	h.activeConnMap.Range(func(key interface{}, value interface{}) bool {
		client := key.(*Client)
		if sender != client.Conn {
			_, err := client.Conn.Write(msg)
			log.Printf("Client in process #%d send %s from port %s", os.Getpid(), string(msg), client.Conn.RemoteAddr().String())
			if err != nil {
				log.Printf("error writing message: %v", err)
			}
		}
		return true
	})
	return nil
}
