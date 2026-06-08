package config

import "flag"

type EvictionPolicy string

const (
	NoEviction     EvictionPolicy = "noeviction"
	AllKeysRandom  EvictionPolicy = "allkeys-random"
	VolatileRandom EvictionPolicy = "volatile-random"
)

// Config contains all runtime configuration required
// to bootstrap and run a Garnet server instance.
type Config struct {
	Host           string
	Port           int
	MaxClients     int
	MaxKeys        int
	EvictionPolicy EvictionPolicy
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

	maxKeys := flag.Int(
		"max-keys",
		DefaultMaxKeys,
		"maximum number of keys to store",
	)

	evictionPolicy := flag.String(
		"eviction-policy",
		string(DefaultEvictionPolicy),
		"eviction policy to use",
	)

	flag.Parse()

	return &Config{
		Host:           *host,
		Port:           *port,
		MaxClients:     *maxClients,
		MaxKeys:        *maxKeys,
		EvictionPolicy: EvictionPolicy(*evictionPolicy),
	}
}
