// Package store implements a small in-memory KV store with vector
// clock semantics. Concurrent writes are NOT resolved — both versions
// are kept, and the client decides how to merge.
package store

import (
	"sync"

	"github.com/thanhtranna/system-design-mastery/examples/vclock-kv/vclock"
)

// Versioned is a single value with its vector clock.
type Versioned struct {
	Value string    `json:"value"`
	Clock vclock.VC `json:"clock"`
}

// Entry is the stored data for one key. May contain multiple concurrent
// versions if writes haven't been reconciled.
type Entry struct {
	Versions []Versioned `json:"versions"`
}

// Store is a thread-safe in-memory KV store.
type Store struct {
	mu     sync.RWMutex
	nodeID string
	data   map[string]*Entry
}

// New returns an empty Store identified as nodeID in vector clocks.
func New(nodeID string) *Store {
	return &Store{
		nodeID: nodeID,
		data:   make(map[string]*Entry),
	}
}

// Put writes a new value. `basedOn` is the set of clocks the client
// observed before deciding to write — used to suppress versions the
// client has reconciled.
func (s *Store) Put(key, value string, basedOn []vclock.VC) Versioned {
	s.mu.Lock()
	defer s.mu.Unlock()

	// New clock = merge of (existing versions OR basedOn) then increment own node.
	existing := s.data[key]
	newClock := vclock.New()
	if existing != nil {
		for _, v := range existing.Versions {
			newClock.Merge(v.Clock)
		}
	}
	for _, c := range basedOn {
		newClock.Merge(c)
	}
	newClock.Increment(s.nodeID)

	newVersion := Versioned{Value: value, Clock: newClock}
	s.applyVersionLocked(key, newVersion)

	return newVersion
}

// ApplyReplicated merges a version received from a peer.
func (s *Store) ApplyReplicated(key string, v Versioned) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applyVersionLocked(key, v)
}

// applyVersionLocked is the core merge: discard versions that are
// strictly dominated, keep concurrent ones.
func (s *Store) applyVersionLocked(key string, v Versioned) {
	entry, ok := s.data[key]
	if !ok {
		s.data[key] = &Entry{Versions: []Versioned{v}}
		return
	}

	kept := make([]Versioned, 0, len(entry.Versions)+1)
	for _, existing := range entry.Versions {
		if existing.Clock.Equal(v.Clock) {
			// Same clock — same write, ignore duplicate.
			return
		}
		if existing.Clock.Dominates(v.Clock) {
			// Existing is newer; ignore the incoming version.
			return
		}
		if v.Clock.Dominates(existing.Clock) {
			// Incoming dominates existing — drop existing.
			continue
		}
		// Concurrent — keep both.
		kept = append(kept, existing)
	}
	kept = append(kept, v)
	entry.Versions = kept
}

// Get returns all current versions of a key.
// If multiple versions exist, the client must reconcile.
func (s *Store) Get(key string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[key]
	if !ok {
		return Entry{}, false
	}
	// Return a copy so caller can't mutate our state.
	out := Entry{Versions: make([]Versioned, len(e.Versions))}
	for i, v := range e.Versions {
		out.Versions[i] = Versioned{Value: v.Value, Clock: v.Clock.Copy()}
	}
	return out, true
}

// All returns all key-versions (for replication / debugging).
func (s *Store) All() map[string]Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Entry, len(s.data))
	for k, e := range s.data {
		copy := Entry{Versions: make([]Versioned, len(e.Versions))}
		for i, v := range e.Versions {
			copy.Versions[i] = Versioned{Value: v.Value, Clock: v.Clock.Copy()}
		}
		out[k] = copy
	}
	return out
}
