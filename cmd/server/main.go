package main

import (
	"log"

	"garnet/internal/config"
	"garnet/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
