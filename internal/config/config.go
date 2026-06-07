package config

import "flag"

// Config contains all runtime configuration required
// to bootstrap and run a Garnet server instance.
type Config struct {
	Host string
	Port int
}

// Load parses command-line flags and returns
// the resulting server configuration.
func Load() *Config {
	host := flag.String(
		"host",
		"0.0.0.0",
		"host address to bind",
	)

	port := flag.Int(
		"port",
		6379,
		"tcp port to listen on",
	)

	flag.Parse()

	return &Config{
		Host: *host,
		Port: *port,
	}
}
