package command

import (
	"garnet/internal/config"
	"garnet/internal/store"
)

func init() {
	// Initialize the store for all tests in the command package
	store.Init(&config.Config{
		MaxKeys:        0,
		EvictionPolicy: config.NoEviction,
	})
}
