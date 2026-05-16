package store_test

import (
	"path/filepath"
	"testing"

	"github.com/sd-book/04-wal/store"
)

func TestSetAndGet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.wal")
	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if err := s.Set("name", "Alice"); err != nil {
		t.Fatal(err)
	}
	if v, ok := s.Get("name"); !ok || v != "Alice" {
		t.Fatalf("got %q %v, want Alice true", v, ok)
	}
}

func TestCrashRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crash.wal")

	// write some data and close (simulating clean shutdown before crash)
	s, _ := store.Open(path)
	s.Set("name", "Alice")
	s.Set("city", "Singapore")
	s.Set("role", "engineer")
	s.Close() // intentional close; next open simulates "recovery after crash"

	// recover
	s2, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	for key, want := range map[string]string{
		"name": "Alice",
		"city": "Singapore",
		"role": "engineer",
	} {
		if got, ok := s2.Get(key); !ok || got != want {
			t.Errorf("key %q: got %q %v, want %q true", key, got, ok, want)
		}
	}
}

func TestDelete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "del.wal")

	s, _ := store.Open(path)
	s.Set("k", "v")
	s.Delete("k")
	s.Close()

	s2, _ := store.Open(path)
	defer s2.Close()
	if _, ok := s2.Get("k"); ok {
		t.Error("expected key to be deleted after recovery")
	}
}
