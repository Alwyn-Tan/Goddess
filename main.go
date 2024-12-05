package main

import (
	"Goddess/server"
	"Goddess/tcp"
	"time"
)

func main() {
	server.Serve(&server.ServerConfig{
		Address:        ":6399",
		MaxConnections: 3,
		Timeout:        10 * time.Second,
	}, tcp.InitEchoHandler())
}
