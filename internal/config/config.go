package config

import "flag"

// Config contains all runtime configuration required
// to bootstrap and run a Garnet server instance.
type Config struct {
	Host       string
	Port       int
	MaxClients int
}

// Load parses command-line flags and returns
// the resulting server configuration.
func Load() *Config {
	host := flag.String(
		"host",
		DefaultHost,
		"host address to bind",
	)

	port := flag.Int(
		"port",
		DefaultPort,
		"tcp port to listen on",
	)

	maxClients := flag.Int(
		"max-clients",
		DefaultMaxClients,
		"maximum concurrent clients for the epoll event loop",
	)

	flag.Parse()

	return &Config{
		Host:       *host,
		Port:       *port,
		MaxClients: *maxClients,
	}
}
