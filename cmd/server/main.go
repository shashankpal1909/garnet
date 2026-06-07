package main

import (
	"log"

	"garnet/internal/config"
	"garnet/internal/server"
)

func main() {
	cfg := config.Load()

	srv := server.New(cfg)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
