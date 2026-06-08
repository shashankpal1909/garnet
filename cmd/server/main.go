package main

import (
	"log"

	"garnet/internal/config"
	"garnet/internal/server"
)

func main() {
	cfg := config.Load()

	var startFunc func() error

	if cfg.Mode == "async" {
		srv := server.NewAsyncServer(cfg)
		startFunc = srv.Start
	} else {
		srv := server.NewSyncServer(cfg)
		startFunc = srv.Start
	}

	if err := startFunc(); err != nil {
		log.Fatal(err)
	}
}
