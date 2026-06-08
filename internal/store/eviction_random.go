package store

import "garnet/internal/logger"

// AllKeysRandomEvictor implements the "allkeys-random" eviction policy.
// It evicts a random key from the entire store regardless of TTL.
type AllKeysRandomEvictor struct{}

// Evict removes one random key from the store.
// Returns 1 if a key was evicted, 0 if the store is empty.
func (e *AllKeysRandomEvictor) Evict(s *Store) int {
	for k := range s.data {
		logger.Logger.Printf("Evicting key '%s' (allkeys-random)", k)
		Delete(k)
		return 1
	}
	return 0
}

// VolatileRandomEvictor implements the "volatile-random" eviction policy.
// It evicts a random key from the set of keys that have an expiration set.
type VolatileRandomEvictor struct{}

// Evict removes one random key with a TTL from the store.
// Returns 1 if a key was evicted, 0 if no keys have an expiration.
func (e *VolatileRandomEvictor) Evict(s *Store) int {
	for k := range s.expires {
		logger.Logger.Printf("Evicting key '%s' (volatile-random)", k)
		Delete(k)
		return 1
	}
	return 0
}
