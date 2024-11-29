package redis

import (
	"Goddess/interface/database"
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type RedisHandler struct {
	activeConnMap sync.Map
	database      database.Database
	closing       atomic.Bool
}

func InitRedisHandler() *RedisHandler {
	var database database.Database
	return &RedisHandler{
		database: database,
	}
}

func (h *RedisHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Load() {
		conn.Close()
		return
	}

	client := StartNewConnection(conn)
	h.activeConnMap.Store(client, struct{}{})

	ch := ParseInputStream(conn)
	for payload := range ch {
		if payload.Error != nil {
		}
	}
}

func (h *RedisHandler) Close() error {
	log.Println("redis handler closing")
	h.closing.Store(true)
	h.activeConnMap.Range(func(key interface{}, value interface{}) bool {
		client := key.(*Connection)
		_ = client.Close()
		return true
	})
	h.database.Close()
	return nil
}
