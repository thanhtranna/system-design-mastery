// Package store provides a WAL-backed key-value store.
// Writes are logged before being applied to the in-memory map.
package store

import (
	"encoding/json"
	"fmt"

	"github.com/sd-book/04-wal/wal"
)

type op struct {
	Op  string `json:"op"`
	Key string `json:"key"`
	Val string `json:"val,omitempty"`
}

// Store is a durable key-value store backed by a WAL.
type Store struct {
	w    *wal.WAL
	data map[string]string
}

// Open opens a store at the given WAL path, replaying any existing log.
func Open(path string) (*Store, error) {
	entries, err := wal.Recover(path)
	if err != nil {
		return nil, fmt.Errorf("recover: %w", err)
	}

	s := &Store{data: make(map[string]string)}
	for _, e := range entries {
		var o op
		if err := json.Unmarshal(e.Data, &o); err != nil {
			return nil, fmt.Errorf("decode entry seq=%d: %w", e.Seq, err)
		}
		s.apply(o)
	}

	w, err := wal.Open(path)
	if err != nil {
		return nil, err
	}
	s.w = w
	return s, nil
}

// Set stores key=value durably.
func (s *Store) Set(key, value string) error {
	return s.log(op{Op: "set", Key: key, Val: value})
}

// Delete removes a key durably.
func (s *Store) Delete(key string) error {
	return s.log(op{Op: "del", Key: key})
}

// Get returns the value for key, and whether it exists.
func (s *Store) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

// Keys returns all keys currently in the store.
func (s *Store) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Close closes the underlying WAL.
func (s *Store) Close() error {
	return s.w.Close()
}

func (s *Store) log(o op) error {
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	if _, err := s.w.Append(b); err != nil {
		return err
	}
	s.apply(o)
	return nil
}

func (s *Store) apply(o op) {
	switch o.Op {
	case "set":
		s.data[o.Key] = o.Val
	case "del":
		delete(s.data, o.Key)
	}
}
