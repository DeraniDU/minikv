// Package store implements MiniKV's in-memory, thread-safe storage.
package store

import "sync"

// Store holds key-value pairs. The mutex protects the map, because
// many goroutines (one per client later) will use it at the same time.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New creates an empty Store. Go has no constructors, so a New function
// is the convention. The map must be created with make before use.
func New() *Store {
	return &Store{data: make(map[string]string)}
}

// Set creates the key or overwrites its value (upsert).
func (s *Store) Set(key, value string) {
	s.mu.Lock() // exclusive lock: nobody else may read or write
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get returns the value and whether the key was found.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock() // shared lock: many readers allowed at once
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

// Delete removes the key and reports whether it existed.
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	delete(s.data, key)
	return ok
}

// Exists reports whether the key is present.
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}

// Len returns the number of keys.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}
