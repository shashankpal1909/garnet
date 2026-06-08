package tests

import (
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestIntegrationSetGet(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=1 to run them.")
	}

	// Connect to our Garnet server (which should be running on localhost:6379 via Docker)
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Fatalf("Could not connect to Garnet server: %v. Make sure 'docker compose up -d' is running!", err)
	}
	defer conn.Close()

	t.Run("Basic SET and GET", func(t *testing.T) {
		// SET a basic key without expiration
		_, err := conn.Do("SET", "integration_key", "garnet_value")
		if err != nil {
			t.Fatalf("Failed to SET: %v", err)
		}

		// GET the key and verify it matches
		val, err := redis.String(conn.Do("GET", "integration_key"))
		if err != nil {
			t.Fatalf("Failed to GET: %v", err)
		}
		if val != "garnet_value" {
			t.Errorf("Expected 'garnet_value', got '%s'", val)
		}
	})

	t.Run("Expiration with EX", func(t *testing.T) {
		// SET a key with a 1-second TTL
		_, err := conn.Do("SET", "temp_key", "temp_value", "EX", 1)
		if err != nil {
			t.Fatalf("Failed to SET with TTL: %v", err)
		}

		// Verify it exists immediately
		val, err := redis.String(conn.Do("GET", "temp_key"))
		if err != nil {
			t.Fatalf("Failed to GET before TTL: %v", err)
		}
		if val != "temp_value" {
			t.Errorf("Expected 'temp_value', got '%s'", val)
		}

		// Verify TTL is reported accurately
		ttl, err := redis.Int(conn.Do("TTL", "temp_key"))
		if err != nil {
			t.Fatalf("Failed to check TTL: %v", err)
		}
		if ttl != 1 && ttl != 0 {
			t.Errorf("Expected TTL of 1 or 0, got %d", ttl)
		}

		// Wait 1.2 seconds for the TTL to expire
		time.Sleep(1200 * time.Millisecond)

		// Verify it has been deleted
		_, err = redis.String(conn.Do("GET", "temp_key"))
		if err != redis.ErrNil {
			t.Errorf("Expected redis.ErrNil error after expiration, got %v", err)
		}

		// Verify TTL returns -2 for deleted keys
		ttlDeleted, err := redis.Int(conn.Do("TTL", "temp_key"))
		if err != nil {
			t.Fatalf("Failed to check TTL: %v", err)
		}
		if ttlDeleted != -2 {
			t.Errorf("Expected TTL of -2, got %d", ttlDeleted)
		}
	})

	t.Run("TTL on key without expiration", func(t *testing.T) {
		// SET a basic key without expiration
		_, err := conn.Do("SET", "forever_key", "val")
		if err != nil {
			t.Fatalf("Failed to SET: %v", err)
		}

		// Verify TTL returns -1
		ttl, err := redis.Int(conn.Do("TTL", "forever_key"))
		if err != nil {
			t.Fatalf("Failed to check TTL: %v", err)
		}
		if ttl != -1 {
			t.Errorf("Expected TTL of -1, got %d", ttl)
		}
	})
}
