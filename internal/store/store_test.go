package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	s := New()
	s.Set("name", "Alex")

	got, ok := s.Get("name")
	if !ok || got != "Alex" {
		t.Fatalf("Get(name) = %q, %v; want %q, true", got, ok, "Alex")
	}
}

func TestGetMissingKey(t *testing.T) {
	s := New()
	if _, ok := s.Get("nope"); ok {
		t.Fatal("expected missing key to return ok=false")
	}
}

func TestSetOverwrites(t *testing.T) {
	s := New()
	s.Set("k", "v1")
	s.Set("k", "v2")

	if got, _ := s.Get("k"); got != "v2" {
		t.Fatalf("got %q, want v2", got)
	}
	if s.Len() != 1 {
		t.Fatalf("Len = %d, want 1", s.Len())
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"existing key", "a", true},
		{"missing key", "zzz", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New()
			s.Set("a", "1")
			if got := s.Delete(tc.key); got != tc.want {
				t.Fatalf("Delete(%q) = %v, want %v", tc.key, got, tc.want)
			}
			if s.Exists(tc.key) {
				t.Fatalf("key %q should not exist after Delete", tc.key)
			}
		})
	}
}

func TestExists(t *testing.T) {
	s := New()
	s.Set("a", "1")
	if !s.Exists("a") {
		t.Fatal("a should exist")
	}
	if s.Exists("b") {
		t.Fatal("b should not exist")
	}
}

// Run with `go test -race` to prove the locking works.
func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			s.Set(key, "value")
			s.Get(key)
			s.Exists(key)
			s.Set("shared", fmt.Sprint(i)) // all goroutines fight over this key
			s.Get("shared")
		}(i)
	}
	wg.Wait()

	if s.Len() != 101 { // 100 unique keys + "shared"
		t.Fatalf("Len = %d, want 101", s.Len())
	}
}
