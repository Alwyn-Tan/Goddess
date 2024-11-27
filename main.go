package main

import (
	"Goddess/tcp"
	"time"
)

func main() {
	tcp.Serve(&tcp.ServerConfig{
		Address:        ":6399",
		MaxConnections: 3,
		Timeout:        10 * time.Second,
	}, tcp.InitEchoHandler())
}
