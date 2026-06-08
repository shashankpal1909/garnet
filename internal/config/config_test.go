package config_test

import (
	"flag"
	"garnet/internal/config"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	oldCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldCommandLine }()

	// Reset flag.CommandLine to avoid panic from redefining flags across tests
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "--host=127.0.0.1", "--port=8080", "--mode=async", "--max-clients=500"}

	cfg := config.Load()

	if cfg.Host != "127.0.0.1" {
		t.Errorf("Expected host 127.0.0.1, got %s", cfg.Host)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Port)
	}

	if cfg.Mode != "async" {
		t.Errorf("Expected mode async, got %s", cfg.Mode)
	}

	if cfg.MaxClients != 500 {
		t.Errorf("Expected max clients 500, got %d", cfg.MaxClients)
	}
}

func TestLoadDefaults(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	oldCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldCommandLine }()

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd"}

	cfg := config.Load()

	if cfg.Host != config.DefaultHost {
		t.Errorf("Expected host %s, got %s", config.DefaultHost, cfg.Host)
	}

	if cfg.Port != config.DefaultPort {
		t.Errorf("Expected port %d, got %d", config.DefaultPort, cfg.Port)
	}

	if cfg.MaxClients != config.DefaultMaxClients {
		t.Errorf("Expected max clients %d, got %d", config.DefaultMaxClients, cfg.MaxClients)
	}
}
