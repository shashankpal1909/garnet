package main

import (
	"log"

	"garnet/internal/config"
	"garnet/internal/server"
	"garnet/internal/store"
)

func main() {
	cfg := config.Load()

	store.Init(cfg)

	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
