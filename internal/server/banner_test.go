package server_test

import (
	"garnet/internal/server"
	"strings"
	"testing"
)

func TestBanner(t *testing.T) {
	if !strings.Contains(server.Banner, "GARNET") && !strings.Contains(server.Banner, "____") {
		t.Errorf("Banner doesn't look like ascii art")
	}
}
