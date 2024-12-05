package server

import (
	"Goddess/database"
	"Goddess/parser"
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Handler struct {
	activeConnMap sync.Map
	database      database.Database
	closing       atomic.Bool
}

func InitHandler() *Handler {
	var db database.Database
	return &Handler{
		database: db,
	}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		conn.Close()
		return
	}

	client := StartNewConnection(conn)
	h.activeConnMap.Store(client, struct{}{})

	ch := parser.ParseInputStream(conn)
	for payload := range ch {
		if payload.Error != nil {
		}
	}
}

func (h *Handler) Close() error {
	log.Println("core handler closing")
	h.closing.Store(true)
	h.activeConnMap.Range(func(key interface{}, value interface{}) bool {
		client := key.(*Connection)
		_ = client.Close()
		return true
	})
	h.database.Close()
	return nil
}
