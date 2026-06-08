package store

import "garnet/internal/config"

// Evictor defines the interface for different memory eviction strategies.
type Evictor interface {
	// Evict attempts to free up space by removing keys from the store.
	// It returns the number of keys successfully evicted.
	Evict(s *Store) int
}

// NoOpEvictor implements the "noeviction" policy, which never evicts any keys.
type NoOpEvictor struct{}

// Evict for NoOpEvictor always returns 0, indicating no keys were removed.
func (*NoOpEvictor) Evict(s *Store) int {
	return 0
}

// NewEvictor creates and returns the appropriate Evictor based on the configured policy.
func NewEvictor(policy config.EvictionPolicy) Evictor {
	switch policy {
	case config.NoEviction:
		return &NoOpEvictor{}
	case config.AllKeysRandom:
		return &AllKeysRandomEvictor{}
	case config.VolatileRandom:
		return &VolatileRandomEvictor{}
	default:
		panic("unknown eviction policy")
	}
}
